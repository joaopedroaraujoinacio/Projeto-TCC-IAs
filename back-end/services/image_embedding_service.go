package services

import (
	"golang_crud/models"
	"golang_crud/repositories"
	"mime/multipart"
	"os"
)

type ImageEmbeddingService struct {
	repo		*repositories.ImageEmbeddingRepository	
	ollamaURL string
}


type OllamaImageRequest struct {
    Model  string   `json:"model"`
    Prompt string   `json:"prompt"`
    Images []string `json:"images"`
    Stream bool     `json:"stream"`
}

type OllamaEmbeddingRequest struct {
    Model  string `json:"model"`
    Prompt string `json:"prompt"`
}

type OllamaEmbeddingRequestResponde struct {
    Embedding []float32 `json:"embedding"`
}

func NewImageEmbeddingService(repo *repositories.ImageEmbeddingRepository, ollamaURL string) *ImageEmbeddingService {
	return &ImageEmbeddingService{
		repo: 			repo,
		ollamaURL: 	ollamaURL,
	}
}

