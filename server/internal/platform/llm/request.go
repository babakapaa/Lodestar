package llm

type Request struct {
	Messages    []Message
	Temperature float64
	MaxTokens   int
	JSON        bool
}
