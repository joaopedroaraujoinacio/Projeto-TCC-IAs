package utils

import (
	"fmt"
	"strings"
)


func VectorToString(vector []float32) string {
	if len(vector) == 0 {
		return "[]"
	} 

	embeddingStr := make([]string, len(vector))
	for i, v := range vector {
		embeddingStr[i] = fmt.Sprintf("%f", v)
	}
	return "[" + strings.Join(embeddingStr, ",") + "]"
}

