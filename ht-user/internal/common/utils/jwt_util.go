package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const secretKey = "]sdf'[8854s1"

func JwtGen(userId int64) string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId":  userId,
			"service": "crm",
			"exp":     time.Now().Add(time.Hour * 2).Unix(),
		})
	signedString, _ := t.SignedString(secretKey)

	return signedString
}
