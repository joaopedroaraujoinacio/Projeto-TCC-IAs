package models

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatRequest struct {
    Message string    `json:"message" binding:"required"`
    History []Message `json:"history"`  // ADD THIS!
    Model   string    `json:"model,omitempty"`
}

type ChatResponse struct {
    Response string `json:"response"`
    Model    string `json:"model"`
}

// Updated for /api/chat endpoint
type OllamaRequest struct {
    Model    string              `json:"model"`
    Messages []map[string]string `json:"messages"` // Changed from Prompt
    Stream   bool                `json:"stream"`
}

// Updated for /api/chat response
type OllamaResponse struct {
    Message struct {
        Content string `json:"content"`
    } `json:"message"` // Changed from Response
    Done bool `json:"done"`
}

type StreamChunk struct {
    Text  string
    Model string
    Done  bool
    Error error
}
