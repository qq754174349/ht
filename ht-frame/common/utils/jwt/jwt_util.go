package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/qq754174349/ht/ht-frame/logger"
)

const defaultSecretKey = "]sdf'[8854s1"

func Gen(claims map[string]interface{}, duration time.Duration, secretKey string) string {
	if secretKey == "" {
		secretKey = defaultSecretKey
	}
	claims["exp"] = time.Now().Add(duration).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	signedString, err := t.SignedString([]byte(secretKey))
	if err != nil {
		log.Error("jwt parse error", err)
	}

	return signedString
}

func Parse(token string, secretKey string) (map[string]interface{}, bool) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if secretKey == "" {
			return []byte(defaultSecretKey), nil
		}
		return secretKey, nil
	}, jwt.WithLeeway(time.Second*10))
	if err != nil {
		return nil, false
	}
	if claims, ok := t.Claims.(jwt.MapClaims); ok {
		return claims, t.Valid
	}
	return nil, false
}
