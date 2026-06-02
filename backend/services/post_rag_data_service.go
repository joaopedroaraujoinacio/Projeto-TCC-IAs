package services

import (
	"fmt"
	"strings"
	"database/sql"
	"go-project/utils"
	"go-project/models"
	"go-project/repositories"
)

const (
	chunkSize    = 1500
	chunkOverlap = 150
)

func CreateDataEmbedding(db *sql.DB, data *models.RagData) (int, error) {
	chunks := chunkText(data.Content, chunkSize, chunkOverlap)

	for i, chunk := range chunks {
		embedding, err := utils.GenerateEmbedding(chunk)
		if err != nil {
			return i, fmt.Errorf("failed to generate embedding for chunk %d: %w", i, err)
		}

		entry := &models.RagData{
			Content:     chunk,
			ContentName: data.ContentName,
			Embedding:   embedding,
		}

		if err := repositories.AddDataToRag(db, entry); err != nil {
			return i, fmt.Errorf("failed to insert chunk %d: %w", i, err)
		}
	}

	return len(chunks), nil
}

// chunkText splits text into overlapping chunks respecting paragraph and sentence boundaries.
func chunkText(text string, size int, overlap int) []string {
	text = strings.TrimSpace(text)
	if len(text) <= size {
		return []string{text}
	}

	// Split into sentences using paragraph and sentence boundaries
	sentences := splitSentences(text)

	var chunks []string
	current := strings.Builder{}

	for _, sentence := range sentences {
		// If adding this sentence would exceed the limit, flush current chunk
		if current.Len()+len(sentence) > size && current.Len() > 0 {
			chunk := strings.TrimSpace(current.String())
			if chunk != "" {
				chunks = append(chunks, chunk)
			}

			// Start next chunk with overlap from the end of current
			tail := current.String()
			if len(tail) > overlap {
				tail = tail[len(tail)-overlap:]
				// Advance to the next word boundary to avoid cutting mid-word
				if idx := strings.Index(tail, " "); idx >= 0 {
					tail = tail[idx+1:]
				}
			}
			current.Reset()
			current.WriteString(tail)
		}

		current.WriteString(sentence)
	}

	if tail := strings.TrimSpace(current.String()); tail != "" {
		chunks = append(chunks, tail)
	}

	return chunks
}

// splitSentences breaks text into sentence-level fragments preserving delimiters.
func splitSentences(text string) []string {
	var sentences []string

	// First split on paragraph breaks
	paragraphs := strings.Split(text, "\n\n")

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		// Split paragraph into sentences
		start := 0
		for i := 0; i < len(para); i++ {
			if i > 0 && (para[i-1] == '.' || para[i-1] == '!' || para[i-1] == '?') && para[i] == ' ' {
				sentence := para[start : i]
				if strings.TrimSpace(sentence) != "" {
					sentences = append(sentences, sentence+" ")
				}
				start = i + 1
			}
		}

		// Remaining text in paragraph (no trailing punctuation)
		if start < len(para) {
			remaining := strings.TrimSpace(para[start:])
			if remaining != "" {
				sentences = append(sentences, remaining+" ")
			}
		}

		// Paragraph separator preserved as whitespace
		sentences = append(sentences, "\n\n")
	}

	return sentences
}
