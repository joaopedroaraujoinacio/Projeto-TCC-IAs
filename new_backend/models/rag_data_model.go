package models


type RagData struct {
	ID int `json:"id"`
	Content string `json:"content" binding:"required"`
	ContentName string `json:"content_name,omitempty"`
	Embedding []float32 `json:"embedding,omitempty"`
}

type EmbeddingRequest struct {
	Model string `json:"model"`
	Prompt string `json:"prompt"`
}

type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

