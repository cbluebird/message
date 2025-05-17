package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"message/config/config"
	"time"
)

func ParseToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(_ *jwt.Token) (interface{}, error) {
		return []byte(config.Config.GetString("jwt.secret")), nil
	})
	if err != nil {
		return -1, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, ok := claims["id"].(int); ok {
			return id, nil
		}
		return -1, jwt.ErrTokenInvalidClaims
	}
	return -1, jwt.ErrTokenInvalidClaims
}

func GenerateToken(id int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := jwt.MapClaims{
		"id":  id,
		"exp": expireTime.Unix(),
		"iat": nowTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := config.Config.GetString("jwt.secret")
	return token.SignedString([]byte(secret))
}
