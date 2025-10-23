package repositories

import (
	"fmt"
	"log"
	"strings"
	"database/sql"
	"go-project/utils"
	"go-project/models"
)


// func RagDataSimilaritySearch(db *sql.DB, queryEmbedding []float32, limit int) ([]models.RagData, error) {
// 	embeddedData := utils.VectorToString(queryEmbedding)
//
// 	log.Printf("embedding string first 100 chars: %s", embeddedData[:100])
//
// 	query := `
//     SELECT id, content,
//            array_cosine_similarity(embedding, CAST(? AS FLOAT[768])) as similarity
//     FROM rag_data 
//     ORDER BY similarity DESC
// 		LIMIT ?
// `
// 	log.Printf("Executing search query with embedding length: %d, limit: %d", len(queryEmbedding), limit)
//
// 	rows, err := db.Query(query, embeddedData, limit)
// 	if err != nil {
// 		log.Printf("search query error: %v", err)
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	var rag_data []models.RagData
// 		for rows.Next() {
// 		var data models.RagData
// 		var similarity float64
//
// 		err := rows.Scan(&data.ID, &data.Content, &similarity)
// 		if err != nil {
// 			log.Printf("Row scan error: %v", err)
// 			continue
// 	}
//
// 	log.Printf("found data ID %d with similarity: %.4f", data.ID, similarity)
// 	rag_data = append(rag_data, data)
//
// 	}	
//
// 	return rag_data, nil
// } 
//

func RagDataSimilaritySearch(db *sql.DB, query string, queryEmbedding []float32, limit int, bm25Weight float64, vectorWeight float64) ([]models.RagData, error) {
	// Tokenize the query for BM25
	tokens := utils.Tokenize(query)
	log.Printf("Search tokens: %v", tokens)
	
	// Get corpus statistics
	var totalDocs int
	var avgDocLength float64
	err := db.QueryRow(`
		SELECT total_docs, avg_content_length 
		FROM corpus_stats 
		WHERE id = 1
	`).Scan(&totalDocs, &avgDocLength)
	if err != nil {
		return nil, fmt.Errorf("failed to get corpus stats: %w", err)
	}
	
	// Get IDF values for query terms
	termIDFs := make(map[string]float64)
	if len(tokens) > 0 {
		placeholders := make([]string, len(tokens))
		args := make([]interface{}, len(tokens))
		for i, token := range tokens {
			placeholders[i] = "?"
			args[i] = token
		}
		
		idfQuery := fmt.Sprintf("SELECT term, doc_freq FROM bm25_stats WHERE term IN (%s)", 
			strings.Join(placeholders, ","))
		
		rows, err := db.Query(idfQuery, args...)
		if err != nil {
			log.Printf("Warning: failed to get term IDFs: %v", err)
		} else {
			defer rows.Close()
			for rows.Next() {
				var term string
				var docFreq int
				if err := rows.Scan(&term, &docFreq); err != nil {
					continue
				}
				termIDFs[term] = utils.CalculateIDF(totalDocs, docFreq)
			}
		}
	}
	
	// Convert query embedding to string for vector search
	embeddedData := utils.VectorToString(queryEmbedding)
	log.Printf("embedding string first 100 chars: %s", embeddedData[:100])
	
	// Perform hybrid search
	searchQuery := `
		SELECT 
			id, 
			content,
			content_name,
			tokens,
			content_length,
			array_cosine_similarity(embedding, CAST(? AS FLOAT[768])) as vector_score
		FROM rag_data
		WHERE tokens IS NOT NULL AND content_length > 0
	`
	
	log.Printf("Executing hybrid search with embedding length: %d, limit: %d", len(queryEmbedding), limit)
	rows, err := db.Query(searchQuery, embeddedData)
	if err != nil {
		log.Printf("search query error: %v", err)
		return nil, err
	}
	defer rows.Close()
	
	type scoredResult struct {
		data        models.RagData
		bm25Score   float64
		vectorScore float64
		hybridScore float64
	}
	
	var results []scoredResult
	
	for rows.Next() {
		var id int
		var content, contentName, tokensStr string
		var docLength int
		var vectorScore float64
		
		err := rows.Scan(&id, &content, &contentName, &tokensStr, &docLength, &vectorScore)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}
		
		// Calculate BM25 score
		docTokens := strings.Fields(tokensStr)
		bm25Score := 0.0
		
		for _, queryTerm := range tokens {
			idf, hasIDF := termIDFs[queryTerm]
			if !hasIDF {
				idf = utils.CalculateIDF(totalDocs, 0)
			}
			
			termFreq := utils.CountTermFrequency(docTokens, queryTerm)
			if termFreq > 0 {
				bm25Score += utils.CalculateBM25Score(termFreq, docLength, avgDocLength, idf)
			}
		}
		
		// Normalize vector score (cosine similarity is already between -1 and 1, convert to 0-1)
		normalizedVectorScore := (vectorScore + 1.0) / 2.0
		
		// Calculate hybrid score
		hybridScore := (bm25Weight * bm25Score) + (vectorWeight * normalizedVectorScore)
		
		log.Printf("Document ID %d - BM25: %.4f, Vector: %.4f, Hybrid: %.4f", 
			id, bm25Score, vectorScore, hybridScore)
		
		results = append(results, scoredResult{
			data: models.RagData{
				ID:          id,
				Content:     content,
				ContentName: contentName,
			},
			bm25Score:   bm25Score,
			vectorScore: vectorScore,
			hybridScore: hybridScore,
		})
	}
	
	// Sort by hybrid score (descending)
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].hybridScore > results[i].hybridScore {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
	
	// Limit results and extract RagData
	var ragData []models.RagData
	for i := 0; i < len(results) && i < limit; i++ {
		ragData = append(ragData, results[i].data)
		log.Printf("Result #%d: ID %d with hybrid score: %.4f", 
			i+1, results[i].data.ID, results[i].hybridScore)
	}
	
	return ragData, nil
}
