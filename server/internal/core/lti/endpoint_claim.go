package lti

type EndpointClaim struct {
	Scopes       []string `json:"scope"`
	LineItemsURL string   `json:"lineitems"`
}
