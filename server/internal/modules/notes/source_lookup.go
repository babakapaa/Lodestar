package notes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	CrossrefURL   = "https://api.crossref.org"
	maxLookupBody = 2 << 20
)

type SourceLookup struct {
	Client   *http.Client
	Crossref string
}

func (l SourceLookup) DOI(ctx context.Context, doi string) (SourceFields, error) {
	doi = NormalizeDOI(doi)
	if !strings.HasPrefix(doi, "10.") || len(doi) > maxSourceText {
		return SourceFields{}, ErrInvalidSource
	}
	body, _, err := l.fetch(ctx, l.Crossref+"/works/"+url.PathEscape(doi), "application/json")
	if err != nil {
		return SourceFields{}, err
	}
	var reply struct {
		Message crossrefWork `json:"message"`
	}
	if err := json.Unmarshal(body, &reply); err != nil {
		return SourceFields{}, fmt.Errorf("decode crossref: %w", ErrLookupFailed)
	}
	return reply.Message.fields(), nil
}

func (l SourceLookup) Page(ctx context.Context, raw string) (SourceFields, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || len(raw) > maxSourceText {
		return SourceFields{}, ErrInvalidSource
	}
	body, contentType, err := l.fetch(ctx, u.String(), "text/html")
	if err != nil {
		return SourceFields{}, err
	}
	fields := SourceFields{Kind: "web", URL: u.String(), AccessedAt: time.Now().UTC().Format(time.DateOnly)}
	if media, _, _ := mime.ParseMediaType(contentType); media == "text/html" || media == "application/xhtml+xml" {
		fields = pageFields(body, fields)
	}
	if fields.Title == "" {
		fields.Title = u.Host
	}
	return fields, nil
}

func (l SourceLookup) fetch(ctx context.Context, target, accept string) ([]byte, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, "", ErrInvalidSource
	}
	req.Header.Set("Accept", accept)
	req.Header.Set("User-Agent", "Lodestar")
	resp, err := l.Client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("fetch source: %w: %w", ErrLookupFailed, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("fetch source: status %d: %w", resp.StatusCode, ErrLookupFailed)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxLookupBody))
	if err != nil {
		return nil, "", fmt.Errorf("read source: %w: %w", ErrLookupFailed, err)
	}
	return body, resp.Header.Get("Content-Type"), nil
}
