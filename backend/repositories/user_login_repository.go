package repositories

import(
	"fmt"
)

func (r *userRepositoryImpl) Save(email, hashedPassword string) (int64, error) {
	query := `INSERT INTO users (email, password, role) VALUES (?, ?, 'user') RETURNING id`
	var id int64
	err := r.db.QueryRow(query, email, hashedPassword).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to save user: %w", err)
	}
	return id, nil
}

func (r *userRepositoryImpl) FindByEmail(email string) (int64, string, string, error) {
	query := `SELECT id, password, role FROM users WHERE email = ?`
	var id int64
	var hashedPassword, role string
	err := r.db.QueryRow(query, email).Scan(&id, &hashedPassword, &role)
	if err != nil {
		return 0, "", "", fmt.Errorf("user not found: %w", err)
	}
	return id, hashedPassword, role, nil
}
