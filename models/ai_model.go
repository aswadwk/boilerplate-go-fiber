package models

type OpenAi struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	// Temperature      int64          `json:"temperature"`
	MaxTokens int64 `json:"max_tokens"`
	// TopP             int64          `json:"top_p"`
	// FrequencyPenalty int64          `json:"frequency_penalty"`
	// PresencePenalty  int64          `json:"presence_penalty"`
	ResponseFormat ResponseFormat `json:"response_format"`
	Stream         bool           `json:"stream"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}
