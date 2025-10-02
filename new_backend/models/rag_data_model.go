package models


type RagData struct {
	ID int `json:"id"`
	Content string `json:"content" binding:"required"`
	ContentName string `json:"content_name,omitempty"`
	Embedding []float32 `json:"embedding,omitempty"`
}

