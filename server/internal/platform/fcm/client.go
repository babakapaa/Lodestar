package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
)

const (
	endpoint = "https://fcm.googleapis.com"
	scope    = "https://www.googleapis.com/auth/firebase.messaging"
)

type Client struct {
	http *http.Client
	url  string
}

func New(serviceAccount []byte) (*Client, error) {
	c, err := parseCredentials(serviceAccount)
	if err != nil {
		return nil, err
	}
	conf := &jwt.Config{Email: c.ClientEmail, PrivateKey: []byte(c.PrivateKey), Scopes: []string{scope}, TokenURL: c.TokenURI}
	base := &http.Client{Timeout: 10 * time.Second}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, base)
	authed := oauth2.NewClient(ctx, conf.TokenSource(ctx))
	authed.Timeout = 10 * time.Second
	return NewWithHTTP(authed, endpoint, c.ProjectID), nil
}

func NewWithHTTP(client *http.Client, baseURL, projectID string) *Client {
	return &Client{http: client, url: baseURL + "/v1/projects/" + projectID + "/messages:send"}
}

func (c *Client) Push(ctx context.Context, token string, m Message) error {
	body, err := json.Marshal(request{Message: message{
		Token:        token,
		Notification: notification{Title: m.Title, Body: m.Body},
		Data:         m.Data,
	}})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("push: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	var failure response
	_ = json.Unmarshal(raw, &failure)
	if resp.StatusCode == http.StatusNotFound || failure.unregistered() {
		return ErrUnregistered
	}
	return fmt.Errorf("push: status %d %s", resp.StatusCode, failure.Error.Status)
}
