package fcm

import (
	"encoding/json"
	"errors"
	"fmt"
)

type credentials struct {
	ProjectID   string `json:"project_id"`
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
	TokenURI    string `json:"token_uri"`
}

func parseCredentials(raw []byte) (credentials, error) {
	var c credentials
	if err := json.Unmarshal(raw, &c); err != nil {
		return credentials{}, fmt.Errorf("read service account: %w", err)
	}
	if c.ProjectID == "" || c.ClientEmail == "" || c.PrivateKey == "" {
		return credentials{}, errors.New("service account needs project_id, client_email and private_key")
	}
	if c.TokenURI == "" {
		c.TokenURI = "https://oauth2.googleapis.com/token"
	}
	return c, nil
}
