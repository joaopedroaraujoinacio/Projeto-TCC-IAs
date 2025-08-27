package services

import (
	"database/sql"
	"golang_crud/models"
	"golang_crud/repositories"
)

type TestService struct {
	repo *repositories.TestRepository
}

func NewTestService(repo *repositories.TestRepository) *TestService {
	return &TestService{repo: repo}
}

func (s *TestService) GetAllTests() ([]models.Test, error) {
	return s.repo.GetAll()
}

func CreateUser(db *sql.DB, user *models.User) error {
	return repositories.CreateUser(db, user)
}

func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	return repositories.GetUserByID(db, id)
}
