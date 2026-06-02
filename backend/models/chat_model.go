package models

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatRequest struct {
		Message 	string 		`json:"message" example:"can you talk a little bit about my grandfather story?"`
    History 	[]Message `json:"history"`  
    Model   	string    `json:"model,omitempty"`
		Provider 	string   	`json:"provider"` 
}

type ChatResponse struct {
    Response string `json:"response"`
    Model    string `json:"model"`
}

type OllamaRequest struct {
    Model    string              `json:"model"`
    Messages []map[string]string `json:"messages"` 
    Stream   bool                `json:"stream"`
}

type OllamaResponse struct {
    Message struct {
        Content string `json:"content"`
    } `json:"message"`
    Done            bool  `json:"done"`
    EvalCount       int   `json:"eval_count"`
    PromptEvalCount int   `json:"prompt_eval_count"`
}

type StreamChunk struct {
    Text         string
    Done         bool
    Error        error
    TokenCount   int
    PromptTokens int
}

type TokenStats struct {
    TokenCount   int `json:"token_count"`
    PromptTokens int `json:"prompt_tokens"`
}

