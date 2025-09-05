
package routes

import (
    "your-project/handlers"
    "github.com/gin-gonic/gin"
)

func SetupImageEmbeddingRoutes(router *gin.Engine, handler *handlers.ImageEmbeddingHandler) {
    api := router.Group("/api/v1")
    {
        images := api.Group("/images")
        {
            images.POST("/embed", handler.EmbedImage)
            images.GET("/search", handler.SearchSimilarImages)
        }
    }
}
```

## 7. Main Application Setup

```go path=/Users/etiennepautet/main.go start_line=1 end_line=15
// Add this to your existing main.go where you initialize services

// Initialize image embedding components
imageEmbeddingRepo := repository.NewImageEmbeddingRepository(db)
imageEmbeddingService := service.NewImageEmbeddingService(imageEmbeddingRepo, "http://ollama:11434")
imageEmbeddingHandler := handlers.NewImageEmbeddingHandler(imageEmbeddingService)

// Setup routes
routes.SetupImageEmbeddingRoutes(router, imageEmbeddingHandler)

// Create uploads directory
if err := os.MkdirAll("./uploads/images", 0755); err != nil {
    log.Fatal("Failed to create uploads directory:", err)
}
```

## 8. Postman Setup

**POST** `http://localhost:8080/api/v1/images/embed`

**Body (form-data):**
- `image` (File): Select your image
- `description` (Text): "A beautiful mountain landscape"
- `source` (Text): "user_upload"
- `tags` (Text): "nature,mountains,landscape"

**Expected Response:**
```json
{
    "id": 1,
    "filename": "1703123456_mountain.jpg",
    "file_path": "./uploads/images/1703123456_mountain.jpg",
    "description": "A beautiful mountain landscape",
    "embedding": [0.1, -0.2, 0.3, ...],
    "metadata": {
        "source": "user_upload",
        "tags": "nature,mountains,landscape"
    },
    "created_at": "2023-12-20T10:30:45Z",
    "message": "Image embedded successfully"
}
