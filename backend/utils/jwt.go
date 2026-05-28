package utils

import (
	"time"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)


const secretKey = "123456789"

func GenerateToken(email string, userId int64, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"userId": userId,
		"role": role,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (int64, string, error) {
	parsedToken, err := jwt.Parse(token,func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, "", err
	}
	
	tokensIsValid := parsedToken.Valid
	if !tokensIsValid {
		return 0, "", errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("invalid token claims")
	}
	userId := int64(claims["userId"].(float64))
	role := claims["role"].(string)

	return userId, role, nil
}

