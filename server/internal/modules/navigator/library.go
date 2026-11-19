package navigator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	OpenAlexURL       = "https://api.openalex.org"
	maxArticles       = 5
	maxArticleAuthors = 3
	maxLibraryQuery   = 200
	maxLibraryBody    = 1 << 20
)

type Library struct {
	Client *http.Client
	URL    string
}

func (l Library) Search(ctx context.Context, query string) ([]Article, error) {
	query = strings.Join(strings.Fields(query), " ")
	if l.Client == nil || query == "" {
		return nil, nil
	}
	if r := []rune(query); len(r) > maxLibraryQuery {
		query = string(r[:maxLibraryQuery])
	}
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	params := url.Values{
		"search":   {query},
		"filter":   {"has_doi:true,is_retracted:false"},
		"per-page": {strconv.Itoa(maxArticles)},
		"select":   {"display_name,doi,publication_year,authorships,primary_location"},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, l.URL+"/works?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("library request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Lodestar")
	resp, err := l.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("library search: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("library search: status %d", resp.StatusCode)
	}
	var reply struct {
		Results []openAlexWork `json:"results"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxLibraryBody)).Decode(&reply); err != nil {
		return nil, fmt.Errorf("decode library: %w", err)
	}
	out := make([]Article, 0, len(reply.Results))
	for _, w := range reply.Results {
		if a, ok := w.article(); ok {
			out = append(out, a)
		}
	}
	return out, nil
}
