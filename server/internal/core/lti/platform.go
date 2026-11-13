package lti

import "time"

type Platform struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Issuer        string    `json:"issuer"`
	ClientID      string    `json:"clientId"`
	AuthURL       string    `json:"authUrl"`
	TokenURL      string    `json:"tokenUrl"`
	JWKSURL       string    `json:"jwksUrl"`
	DeploymentIDs []string  `json:"deploymentIds"`
	MatchEmail    bool      `json:"matchEmail"`
	CreatedAt     time.Time `json:"createdAt"`
}

func (p Platform) Deploys(id string) bool {
	for _, d := range p.DeploymentIDs {
		if d == id {
			return true
		}
	}
	return false
}
