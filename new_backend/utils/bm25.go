package utils

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"regexp"
	"strings"
)

const (
	K1 = 1.2  
	B  = 0.75 
)

var stopwords = map[string]bool{
	"a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
	"be": true, "by": true, "for": true, "from": true, "has": true, "he": true,
	"in": true, "is": true, "it": true, "its": true, "of": true, "on": true,
	"that": true, "the": true, "to": true, "was": true, "will": true, "with": true,
	"or": true, "but": true, "not": true, "this": true, "they": true, "been": true,
}

func Tokenize(text string) []string {
	text = strings.ToLower(text)
	
	reg := regexp.MustCompile("[^a-z0-9\\s]+")
	text = reg.ReplaceAllString(text, " ")
	
	words := strings.Fields(text)
	
	var tokens []string
	for _, word := range words {
		if len(word) >= 2 && !stopwords[word] {
			tokens = append(tokens, word)
		}
	}
	
	return tokens
}

func TokensToString(tokens []string) string {
	return strings.Join(tokens, " ")
}

func CountTermFrequency(tokens []string, term string) int {
	count := 0
	for _, t := range tokens {
		if t == term {
			count++
		}
	}
	return count
}

func CalculateIDF(totalDocs int, docFreq int) float64 {
	if docFreq == 0 {
		return 0.0
	}
	numerator := float64(totalDocs) - float64(docFreq) + 0.5
	denominator := float64(docFreq) + 0.5
	
	if numerator <= 0 || denominator <= 0 {
		return 0.0
	}
	
	return math.Log(numerator / denominator)
}

func CalculateBM25Score(termFreq int, docLength int, avgDocLength float64, idf float64) float64 {
	if avgDocLength == 0 {
		return 0.0
	}
	
	numerator := float64(termFreq) * (K1 + 1)
	denominator := float64(termFreq) + K1*(1-B+B*float64(docLength)/avgDocLength)
	
	if denominator == 0 {
		return 0.0
	}
	
	return idf * (numerator / denominator)
}

func UpdateBM25Stats(db *sql.DB) error {
	log.Println("Updating BM25 statistics...")
	
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	if _, err := tx.Exec("DELETE FROM bm25_stats"); err != nil {
		return fmt.Errorf("failed to clear bm25_stats: %w", err)
	}
	
	var totalDocs int
	var avgDocLength float64
	err = tx.QueryRow(`
		SELECT 
			COUNT(*) as total_docs,
			COALESCE(AVG(doc_length), 0.0) as avg_doc_length
		FROM rag_data
		WHERE tokens IS NOT NULL AND doc_length > 0
	`).Scan(&totalDocs, &avgDocLength)
	if err != nil {
		return fmt.Errorf("failed to calculate corpus stats: %w", err)
	}
	
	_, err = tx.Exec(`
		UPDATE corpus_stats 
		SET total_docs = ?, 
		    avg_doc_length = ?,
		    last_updated = CURRENT_TIMESTAMP
		WHERE id = 1
	`, totalDocs, avgDocLength)
	if err != nil {
		return fmt.Errorf("failed to update corpus_stats: %w", err)
	}
	
	rows, err := tx.Query("SELECT tokens FROM rag_data WHERE tokens IS NOT NULL")
	if err != nil {
		return fmt.Errorf("failed to query documents: %w", err)
	}
	defer rows.Close()
	
	termDocFreq := make(map[string]int)    
	termTotalFreq := make(map[string]int) 
	
	for rows.Next() {
		var tokensStr string
		if err := rows.Scan(&tokensStr); err != nil {
			continue
		}
		
		tokens := strings.Fields(tokensStr)
		seenInDoc := make(map[string]bool)
		
		for _, token := range tokens {
			termTotalFreq[token]++
			if !seenInDoc[token] {
				termDocFreq[token]++
				seenInDoc[token] = true
			}
		}
	}
	
	if len(termDocFreq) > 0 {
		values := make([]string, 0, len(termDocFreq))
		args := make([]interface{}, 0, len(termDocFreq)*3)
		
		for term, docFreq := range termDocFreq {
			values = append(values, "(?, ?, ?)")
			args = append(args, term, docFreq, termTotalFreq[term])
		}
		
		query := fmt.Sprintf("INSERT INTO bm25_stats (term, doc_freq, total_freq) VALUES %s", 
			strings.Join(values, ","))
		
		_, err = tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("failed to insert term stats: %w", err)
		}
	}
	
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	log.Printf(" BM25 stats updated: %d documents, %.2f avg length, %d unique terms",
		totalDocs, avgDocLength, len(termDocFreq))
	
	return nil
}

