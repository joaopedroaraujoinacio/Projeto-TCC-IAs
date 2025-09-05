package services

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "time"
    "your-project/models"
    "your-project/repository"
)

type ImageEmbeddingService struct {
    repo      *repository.ImageEmbeddingRepository
    ollamaURL string
}

type OllamaImageRequest struct {
    Model  string   `json:"model"`
    Prompt string   `json:"prompt"`
    Images []string `json:"images"`
    Stream bool     `json:"stream"`
}

type OllamaImageResponse struct {
    Response string `json:"response"`
    Done     bool   `json:"done"`
}

type OllamaEmbeddingRequest struct {
    Model  string `json:"model"`
    Prompt string `json:"prompt"`
}

type OllamaEmbeddingResponse struct {
    Embedding []float32 `json:"embedding"`
}

func NewImageEmbeddingService(repo *repository.ImageEmbeddingRepository, ollamaURL string) *ImageEmbeddingService {
    return &ImageEmbeddingService{
        repo:      repo,
        ollamaURL: ollamaURL,
    }
}

func (s *ImageEmbeddingService) EmbedImage(file multipart.File, header *multipart.FileHeader, request *models.ImageEmbedRequest) (*models.ImageEmbedResponse, error) {
    // Save file to disk
    uploadDir := "./uploads/images"
    if err := os.MkdirAll(uploadDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create upload directory: %v", err)
    }

    filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
    filePath := filepath.Join(uploadDir, filename)
    
    dst, err := os.Create(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to create file: %v", err)
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        return nil, fmt.Errorf("failed to save file: %v", err)
    }

    // Convert image to base64
    imageData, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read image file: %v", err)
    }
    base64Image := base64.StdEncoding.EncodeToString(imageData)

    // Generate embedding using LLaVA
    embedding, err := s.generateImageEmbedding(base64Image, request.Description)
    if err != nil {
        return nil, fmt.Errorf("failed to generate embedding: %v", err)
    }

    // Create database record
    imageEmbedding := &models.ImageEmbedding{
        Filename:    filename,
        FilePath:    filePath,
        Description: &request.Description,
        Embedding:   embedding,
        Metadata:    request.Metadata,
    }

    if err := s.repo.Create(imageEmbedding); err != nil {
        return nil, fmt.Errorf("failed to save to database: %v", err)
    }

    return &models.ImageEmbedResponse{
        ImageEmbedding: *imageEmbedding,
        Message:        "Image embedded successfully",
    }, nil
}

func (s *ImageEmbeddingService) generateImageEmbedding(base64Image, description string) ([]float32, error) {
    prompt := "Generate a detailed description of this image for semantic search purposes."
    if description != "" {
        prompt = fmt.Sprintf("Analyze this image with context: %s. Provide a detailed description for semantic search.", description)
    }

    ollamaReq := OllamaImageRequest{
        Model:  "llava:7b",
        Prompt: prompt,
        Images: []string{base64Image},
        Stream: false,
    }

    reqBody, err := json.Marshal(ollamaReq)
    if err != nil {
        return nil, err
    }

    resp, err := http.Post(
        fmt.Sprintf("%s/api/generate", s.ollamaURL),
        "application/json",
        bytes.NewBuffer(reqBody),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var ollamaResp OllamaImageResponse
    if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
        return nil, err
    }

    // Now generate embedding from the description
    return s.generateTextEmbedding(ollamaResp.Response)
}

func (s *ImageEmbeddingService) generateTextEmbedding(text string) ([]float32, error) {
    embedReq := OllamaEmbeddingRequest{
        Model:  "nomic-embed-text:latest",
        Prompt: text,
    }

    reqBody, err := json.Marshal(embedReq)
    if err != nil {
        return nil, err
    }

    resp, err := http.Post(
        fmt.Sprintf("%s/api/embeddings", s.ollamaURL),
        "application/json",
        bytes.NewBuffer(reqBody),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var embedResp OllamaEmbeddingResponse
    if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
        return nil, err
    }

    return embedResp.Embedding, nil
}

func (s *ImageEmbeddingService) SearchSimilarImages(queryEmbedding []float32, limit int) ([]models.ImageEmbedding, error) {
    return s.repo.FindSimilar(queryEmbedding, limit)
}
