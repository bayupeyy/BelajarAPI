package middlewares

import (
	"21-api/config"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(hp string) (string, error) {
	var data = jwt.MapClaims{}
	// custom data
	data["hp"] = hp
	// mandatory data
	data["iat"] = time.Now().Unix()
	data["exp"] = time.Now().Add(time.Hour * 3).Unix()

	var proccessToken = jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	result, err := proccessToken.SignedString([]byte(config.JWTSECRET))

	if err != nil {
		defer func() {
			if err := recover(); err != nil {
				log.Println("error jwt creation", err)
			}
		}()
		return "", errors.New("Terjadi masalah")
	}

	return result, nil
}

func DecodeToken(token *jwt.Token) string {
	var result string
	var claim = token.Claims.(jwt.MapClaims)

	if val, found := claim["hp"]; found {
		result = val.(string)
	}

	return result
}
