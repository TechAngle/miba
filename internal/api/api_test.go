package api

import (
	"testing"
)

func TestParseCredentials(t *testing.T) {
	client, err := NewAPIClient()
	if err != nil {
		t.Fatal(err)
	}

	client.UpdateInformation()
}
