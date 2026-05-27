package models

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}
