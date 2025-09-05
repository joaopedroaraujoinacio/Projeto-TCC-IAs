package handlers

import (
    "net/http"
    "strconv"
    "your-project/models"
    "your-project/service"
    "github.com/gin-gonic/gin"
)

type ImageEmbeddingHandler struct {
    service *service.ImageEmbeddingService
}

func NewImageEmbeddingHandler(service *service.ImageEmbeddingService) *ImageEmbeddingHandler {
    return &ImageEmbeddingHandler{service: service}
}

func (h *ImageEmbeddingHandler) EmbedImage(c *gin.Context) {
    // Parse multipart form
    file, header, err := c.Request.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
        return
    }
    defer file.Close()

    // Validate file type
    if !isValidImageType(header.Header.Get("Content-Type")) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image type. Only JPEG, PNG, and WebP are supported"})
        return
    }

    // Parse additional form data
    request := &models.ImageEmbedRequest{
        Description: c.PostForm("description"),
        Metadata:    make(map[string]interface{}),
    }

    // Parse metadata if provided
    if source := c.PostForm("source"); source != "" {
        request.Metadata["source"] = source
    }
    if tags := c.PostForm("tags"); tags != "" {
        request.Metadata["tags"] = tags
    }

    // Process the image
    response, err := h.service.EmbedImage(file, header, request)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, response)
}

func (h *ImageEmbeddingHandler) SearchSimilarImages(c *gin.Context) {
    limitStr := c.DefaultQuery("limit", "10")
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 {
        limit = 10
    }

    // This would typically take a text query and convert it to embedding
    // For now, this is a placeholder
    c.JSON(http.StatusOK, gin.H{"message": "Search endpoint - implementation depends on your search requirements"})
}

func isValidImageType(contentType string) bool {
    validTypes := []string{
        "image/jpeg",
        "image/jpg", 
        "image/png",
        "image/webp",
    }
    
    for _, validType := range validTypes {
        if contentType == validType {
            return true
        }
    }
    return false
}
