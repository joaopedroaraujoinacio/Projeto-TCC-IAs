package models

type Document struct {
    ID        int       `json:"id"`
    Content   string    `json:"content" binding:"required"`
    MediaType *string    `json:"media_type"`
    FileName  *string    `json:"file_name,omitempty"`
    Embedding []float32 `json:"embedding,omitempty"`
    CreatedAt string    `json:"created_at"`
}
