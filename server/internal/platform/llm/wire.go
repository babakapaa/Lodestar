package llm

type chatBody struct {
	Model          string          `json:"model"`
	Messages       []Message       `json:"messages"`
	Temperature    float64         `json:"temperature"`
	MaxTokens      int             `json:"max_tokens,omitempty"`
	Stream         bool            `json:"stream"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type chatReply struct {
	Choices []struct {
		Message Message `json:"message"`
		Delta   Message `json:"delta"`
	} `json:"choices"`
}

func (r chatReply) text(stream bool) string {
	if len(r.Choices) == 0 {
		return ""
	}
	if stream {
		return r.Choices[0].Delta.Content
	}
	return r.Choices[0].Message.Content
}
