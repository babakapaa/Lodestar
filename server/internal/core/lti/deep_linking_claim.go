package lti

type DeepLinkingClaim struct {
	ReturnURL      string   `json:"deep_link_return_url"`
	AcceptTypes    []string `json:"accept_types"`
	AcceptMultiple *bool    `json:"accept_multiple"`
	Data           string   `json:"data"`
}

func (c DeepLinkingClaim) acceptsLinks() bool {
	for _, t := range c.AcceptTypes {
		if t == "ltiResourceLink" {
			return true
		}
	}
	return false
}
