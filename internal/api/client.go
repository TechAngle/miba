package api

import (
	"fmt"
	"io"
	"miba/internal/config"
	"miba/internal/domain"
	"net/http"
	"net/url"
	"regexp"
	"time"

	cookiejar "github.com/juju/persistent-cookiejar"
)

// MiAPIClient is a main orchestrator of requests for API.
type MiAPIClient struct {
	// Stok is a key used by Mi Router to validate request session.
	Stok string `json:"stok"`

	baseURL     *url.URL
	client      *http.Client
	credentials *domain.MiAPICredentials
}

// NewAPIClient creates new API client with prepared client
func NewAPIClient() (*MiAPIClient, error) {
	client, err := newHttpClient()
	if err != nil {
		return nil, err
	}

	baseURL, _ := url.Parse(domain.BaseURL)

	return &MiAPIClient{
		baseURL:     baseURL,
		client:      client,
		credentials: &domain.MiAPICredentials{},
	}, nil
}

// PingRouter checks if router is available. Returns error if not.
func (c *MiAPIClient) PingRouter() error {
	return c.routerAvailable()
}

// Login generates nonce and encrypts password with it, after that logins and updates
// client stok token.
func (c *MiAPIClient) Login(password string) (stok string, err error) {
	nonce := generateNonce("0", c.credentials.DeviceId)
	encodedPassword := encryptPassword(password, c.credentials.Key, nonce)

	params := loginParams(encodedPassword, nonce)
	data, err := c.loginRequest(params)
	if err != nil {
		return "", err
	}

	/*
		 I hate xiaomi for their shitcode tbh.
			Who, the hell, returns 200 OK but 401 in structure? Are you serious?
	*/
	if data.Code == http.StatusUnauthorized {
		return "", domain.ErrUnauthorized
	}

	if data.Code != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", data.Code)
	}

	// if token is empty showing debug info
	if data.Token == "" {
		return "", domain.ErrEmptyToken
	}

	c.Stok = data.Token

	return data.Token, nil
}

// UpdateInformation looks up for credentials from luci web page using regexes.
// Should be executed BEFORE log in.
func (c *MiAPIClient) UpdateInformation() error {
	content, err := c.luciWebContent()
	if err != nil {
		return err
	}

	systemContext, err := c.parseSystemContext(content)
	if err != nil {
		return err
	}

	c.credentials.IV = string(systemContext.IV)
	c.credentials.Key = string(systemContext.Key)
	c.credentials.DeviceId = string(systemContext.DeviceID)

	return nil
}

// loginParams creates url.Values structure with all needed fields.
func loginParams(password string, nonce string) url.Values {
	params := url.Values{}
	params.Add("username", "admin")
	params.Add("password", password)
	params.Add("logtype", "2")
	params.Add("nonce", nonce)

	return params
}

// parseSystemContext parses content and extracts fields for system context.
//
// Returns an error if context is invalid.
func (c *MiAPIClient) parseSystemContext(content []byte) (*domain.SystemContext, error) {
	var context domain.SystemContext

	fields := []struct {
		dest *string
		re   *regexp.Regexp
	}{
		{dest: &context.DeviceID, re: domain.DeviceIDRegex},
		{dest: &context.HardwareID, re: domain.HardwareRegex},
		{dest: &context.RomVersion, re: domain.ROMVersionRegex},
		{dest: &context.MAC, re: domain.MacRegex},
		{dest: &context.Key, re: domain.KeyRegex},
		{dest: &context.IV, re: domain.IvRegex},
	}

	for _, field := range fields {
		*field.dest = string(findMatch(field.re, content))
	}

	if !context.Valid() {
		return nil, domain.ErrInvalidSystemContext
	}

	return &context, nil
}

// findMatch uses FindSubmatch for regexp and returns nil if length of matches
// less than 2. Otherwise, if match was found it will be returned.
func findMatch(r *regexp.Regexp, v []byte) []byte {
	match := r.FindSubmatch(v)
	if len(match) < 2 {
		return nil
	}

	return match[1]
}

// luciWebContent retrieves all page content from Mi Router luci/web page.
func (c *MiAPIClient) luciWebContent() ([]byte, error) {
	endpoint := c.buildUrl(domain.WebEndpoint, nil)

	res, err := c.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("mi endpoint err: %w", err)
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read mi body err: %w", err)
	}

	return content, nil
}

// routerAvailable checks Mi Router endpoint using HEAD request to it.
//
// Returns error if something went wrong or has invalid status code.
func (c *MiAPIClient) routerAvailable() error {
	res, err := c.client.Head(c.baseURL.String())
	if err != nil {
		return fmt.Errorf("router head request err: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status not ok (%d)", res.StatusCode)
	}

	return nil
}

func newHttpClient() (*http.Client, error) {
	cookieJar, err := cookiejar.New(&cookiejar.Options{
		Filename: config.CookieJarBase,
	})
	if err != nil {
		return nil, fmt.Errorf("new cookie jar err: %w", err)
	}

	client := http.Client{
		Timeout: 15 * time.Second,
		Jar:     cookieJar,
	}

	return &client, nil
}
