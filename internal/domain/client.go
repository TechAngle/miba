package domain

import (
	"slices"
)

// SystemContext contains all machine information including key and iv for passwords
type SystemContext struct {
	DeviceID   string
	HardwareID string
	RomVersion string
	MAC        string
	Key        string
	IV         string
}

// Valid checks if context contains any empty field (usually it should not be).
func (c SystemContext) Valid() bool {
	return !anyEmpty(
		c.DeviceID,
		c.HardwareID,
		c.RomVersion,
		c.MAC,
		c.Key,
		c.IV,
	)
}

// MiAPICredentials contains all credentials used for requests authorization.
type MiAPICredentials struct {
	// DeviceId is device MAC
	DeviceId string

	// Key is a randomly generated key.
	Key string

	// IV is a key used for 'newPwd' field.
	IV string
}

// Valid checks if none of fields is empty.
func (c MiAPICredentials) Valid() bool {
	return !anyEmpty(c.IV, c.Key, c.DeviceId)
}

// anyEmpty checks if elements have at least one empty string
func anyEmpty(elements ...string) bool {
	return slices.Contains(elements, "")
}
