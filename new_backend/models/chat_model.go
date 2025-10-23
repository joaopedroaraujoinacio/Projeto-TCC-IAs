package models

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatRequest struct {
    Message string    `json:"message" binding:"required"`
    History []Message `json:"history"`  
    Model   string    `json:"model,omitempty"`
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
    Done bool `json:"done"`
}

type StreamChunk struct {
    Text  string
    Model string
    Done  bool
    Error error
}

