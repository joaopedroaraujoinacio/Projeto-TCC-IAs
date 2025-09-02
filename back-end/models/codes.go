package models

type Codes struct {
	ID        					int       	`json:"id"`
	Content       			string      `json:"content" binding:"required"`
	ProgramingLanguage 	string			`json:"programming_language"`
	FileName  					string      `json:"file_name,omitempty"`
	Embedding        		[]float32		`json:"embedding,omitempty"`
	CreatedAt        		string      `json:"created_at"`
}
