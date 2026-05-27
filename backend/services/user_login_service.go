package services

import (
	"fmt"
	"go-project/models"
	"go-project/utils"
)

func (s *UserService) Signup(user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	id, err := s.repo.Save(user.Email, hashedPassword)
	if err != nil {
		return err
	}
	user.ID = id
	user.Role = "user"
	return nil
}

func (s *UserService) Login(email, password string) (*models.User, error) {
	id, hashedPassword, role, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	if !utils.CheckPasswordHash(password, hashedPassword) {
		return nil, fmt.Errorf("invalid credentials")
	}
	return &models.User{ID: id, Email: email, Role: role}, nil
}
