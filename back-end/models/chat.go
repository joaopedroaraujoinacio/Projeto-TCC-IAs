package models


type ChatRequest struct {
	Message string `json:"message" binding:"required"`
	Model   string `json:"model,omitempty"`
}

type ChatResponse struct {
	Response string `json:"response"`
	Model    string `json:"model"`
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

