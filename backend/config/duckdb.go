package config 

import (
	"os"
	"fmt"
	"log"
	"path/filepath"
	"database/sql"
	_ "github.com/marcboeker/go-duckdb"
)


func ConnectDB(databaseURL string) (*sql.DB, error) {
	if err := ensureDBDirectory(databaseURL)
	err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("duckdb", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); 
	err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err) 
	}


	log.Println("Installing vss for vector operations...")
	if _, err := db.Exec("INSTALL vss");
	err != nil {
		log.Println("vss might already be installed: %w", err)
	}

	if _, err :=db.Exec("LOAD vss");
	err != nil {
		return nil, fmt.Errorf("failed to load vss extension: %w", err)
	}

	if err := initializeSchema(db);
	err != nil {
		return nil, fmt.Errorf("failed to initializeSchema: %w", err)
	}

	log.Println("successfully connected to DuckDB with vector support")
	return db, nil
}

func ensureDBDirectory(databaseURL string) error {
	if databaseURL != ":memory:" {
		dir := filepath.Dir(databaseURL)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func initializeSchema(db *sql.DB) error {
	schema := `
	CREATE SEQUENCE IF NOT EXISTS sequence_users START 1;
	CREATE TABLE IF NOT EXISTS users (
		id       INTEGER PRIMARY KEY DEFAULT nextval('sequence_users'),
		email    VARCHAR(255) NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role     VARCHAR(50) NOT NULL DEFAULT 'user'
	);

	CREATE SEQUENCE IF NOT EXISTS sequence_rag_data START 1;
	CREATE TABLE IF NOT EXISTS rag_data (
		id             INTEGER PRIMARY KEY DEFAULT nextval('sequence_rag_data'),
		user_id        INTEGER NOT NULL REFERENCES users(id),
		content        TEXT NOT NULL,
		content_name   VARCHAR(100),
		embedding      FLOAT[768],
		tokens         TEXT,
		content_length INTEGER
	);

	CREATE TABLE IF NOT EXISTS bm25_stats (
		user_id    INTEGER NOT NULL REFERENCES users(id),
		term       TEXT NOT NULL,
		doc_freq   INTEGER,
		total_freq INTEGER,
		PRIMARY KEY (user_id, term)
	);

	CREATE TABLE IF NOT EXISTS corpus_stats (
		user_id            INTEGER PRIMARY KEY REFERENCES users(id),
		total_docs         INTEGER,
		avg_content_length FLOAT
	);

	CREATE INDEX IF NOT EXISTS idx_rag_data_user ON rag_data(user_id);
	CREATE INDEX IF NOT EXISTS idx_rag_data_tokens ON rag_data(tokens);
	`
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	log.Println("database schema initialized successfully")
	return nil
}

