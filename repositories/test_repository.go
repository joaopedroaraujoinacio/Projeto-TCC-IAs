package repositories

import (
	"database/sql"
	"golang_crud/models"
)

type TestRepository struct {
	db *sql.DB
}

func NewTestRepository(db *sql.DB) *TestRepository {
	return &TestRepository{db: db}
}

func (r *TestRepository) GetAll() ([]models.Test, error) {
	query := `SELECT id, message, created_at FROM tests ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []models.Test
	for rows.Next() {
		var test models.Test
		if err := rows.Scan(&test.ID, &test.Message, &test.CreatedAt); err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}
	
	return tests, nil
}

func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	var user models.User
  query := `SELECT id, name, email FROM users WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(db *sql.DB, user *models.User) error {
  query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
	err := db.QueryRow(query, user.Name, user.Email).Scan(&user.ID)
	return err
} 
