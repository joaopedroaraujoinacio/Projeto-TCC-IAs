package models

type StreamOptions struct {
    IncludeUsage bool `json:"include_usage"`
}

type OpenAIRequest struct {
    Model         string              `json:"model"`
    Messages      []map[string]string `json:"messages"`
    Stream        bool                `json:"stream"`
    StreamOptions *StreamOptions      `json:"stream_options,omitempty"`
}

type OpenAIStreamChunk struct {
    Choices []struct {
        Delta struct {
            Content string `json:"content"`
        } `json:"delta"`
        FinishReason *string `json:"finish_reason"`
    } `json:"choices"`
    Usage *struct {
        PromptTokens     int `json:"prompt_tokens"`
        CompletionTokens int `json:"completion_tokens"`
        TotalTokens      int `json:"total_tokens"`
    } `json:"usage"`
}
