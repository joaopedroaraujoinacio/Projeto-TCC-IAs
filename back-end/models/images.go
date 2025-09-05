package models


type ImageEmbeding struct {
    ID          int                    `json:"id"`
    Filename    string                 `json:"filename"`
    FilePath    string                 `json:"file_path"`
    Description *string                `json:"description"`
    Embedding   []float32              `json:"embedding,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
    CreatedAt   string                 `json:"created_at"`
}

type ImageEmbedRequest struct {
    Description string                 `json:"description"`
    Metadata    map[string]interface{} `json:"metadata"`
}

type ImageEmbedResponse struct {
		ImageEmbeding
    Message string `json:"messsage"`
}

