package lti

import (
	"slices"
	"strings"
)

type PlatformInput struct {
	Name          string   `json:"name"`
	Issuer        string   `json:"issuer"`
	ClientID      string   `json:"clientId"`
	AuthURL       string   `json:"authUrl"`
	TokenURL      string   `json:"tokenUrl"`
	JWKSURL       string   `json:"jwksUrl"`
	DeploymentIDs []string `json:"deploymentIds"`
	MatchEmail    bool     `json:"matchEmail"`
}

func (in PlatformInput) Normalized() (PlatformInput, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Issuer = strings.TrimSpace(in.Issuer)
	in.ClientID = strings.TrimSpace(in.ClientID)
	in.AuthURL = strings.TrimSpace(in.AuthURL)
	in.TokenURL = strings.TrimSpace(in.TokenURL)
	in.JWKSURL = strings.TrimSpace(in.JWKSURL)
	ids := make([]string, 0, len(in.DeploymentIDs))
	for _, id := range in.DeploymentIDs {
		if id = strings.TrimSpace(id); id != "" {
			ids = append(ids, id)
		}
	}
	slices.Sort(ids)
	in.DeploymentIDs = slices.Compact(ids)
	if len(in.DeploymentIDs) == 0 {
		return PlatformInput{}, ErrInvalidPlatform
	}
	return in, nil
}
