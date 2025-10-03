package utils

import (
	"fmt"
	"strings"
	"strconv"
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

func ParseVectorString(vectorSTR string) []float32 {
	vectorSTR = strings.Trim(vectorSTR, "[]")
	if vectorSTR == "" {
		return []float32{}
	}

	parts := strings.Split(vectorSTR, ",")
	embedding := make([]float32, len(parts))
	
	for i, part := range  parts {
		if val, err := strconv.ParseFloat(strings.TrimSpace(part), 32); err == nil {
			embedding[i] = float32(val)
		}
	}

	return embedding
}

