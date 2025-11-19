package services

import (
	// "errors"
	"database/sql"
	"go-project/utils"
)


type User struct {
	ID int64
	Email string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save(db *sql.DB) error {
	query := "INSERT INTO users(email, password) VALUES (?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPassword)
	if err != nil {
		return err
	}

	userID, err := result.LastInsertId()

	u.ID = userID
	return err
}

// func (u *User) ValidateCredentials() error {
// 	query := "SELECT id, password FROM users WHERE email = ?"
// 	row := db.DB.QueryRow(query, u.Email)
//
// 	var retreivedPassword string
// 	err := row.Scan(&u.ID, &retreivedPassword)
// 	if err != nil {
// 		return errors.New("error when retrieving user data")
// 	}
//
// 	passwordIsValid := utils.CheckPasswordHash(u.Password, retreivedPassword)
//
// 	if !passwordIsValid {
// 		return errors.New("passwords dont match")
// 	}
//
// 	return nil
// }
//
