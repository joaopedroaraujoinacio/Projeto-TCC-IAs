package utils

import (
	"fmt"
	"strconv"
	"strings"
)


func VectorToString(embedding []float32) string {
	strValues := make([]string, len(embedding))
	for i, v := range embedding {
		strValues[i] = fmt.Sprintf("%f", v)
	}
	return "[" + strings.Join(strValues, ",") + "]"
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
