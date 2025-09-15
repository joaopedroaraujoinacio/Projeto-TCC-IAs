CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    media_type VARCHAR(100),
    file_name VARCHAR(255),
    embedding vector(768),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS documents_embedding_idx 
ON documents USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 1);

CREATE OR REPLACE FUNCTION search_similar_documents(
    query_embedding vector(768),
    match_count int DEFAULT 5
) RETURNS TABLE (
    id INT,
    content TEXT,
    similarity FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        d.id,
        d.content,
        1 - (d.embedding <=> query_embedding) AS similarity
    FROM documents d
    WHERE d.embedding IS NOT NULL
    ORDER BY d.embedding <=> query_embedding
    LIMIT match_count;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE IF NOT EXISTS codes (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    programming_language VARCHAR(100),
    file_name VARCHAR(255),
    embedding vector(1024),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS codes_embedding_idx 
ON codes USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 1);

CREATE OR REPLACE FUNCTION search_similar_code_documents(
    query_embedding vector(1024),
    match_count int DEFAULT 5
) RETURNS TABLE (
    id INT,
    content TEXT,
    similarity FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        d.id,
        d.content,
        1 - (d.embedding <=> query_embedding) AS similarity
    FROM documents d
    WHERE d.embedding IS NOT NULL
    ORDER BY d.embedding <=> query_embedding
    LIMIT match_count;
END;
$$ LANGUAGE plpgsql;

