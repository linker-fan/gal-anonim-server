package auth

import (
	"errors"
	"fmt"
	"linker-fan/gal-anonim-server/server/config"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.StandardClaims
}

var c *config.Config

func init() {
	conf, err := config.NewConfig("./config.yml")
	if err != nil {
		log.Fatal(err)
	}

	c = conf
}

func GenerateJWTToken(username string, id int, isAdmin bool) (string, *time.Time, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(c.Jwt.ExpTime))
	claims := &JwtClaims{
		ID:       id,
		Username: username,
		IsAdmin:  isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    c.Jwt.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.Jwt.TokenSecret))
	if err != nil {
		log.Println(err)
		return "", nil, err
	}

	return tokenString, &expirationTime, nil
}

func IsTokenValid(tokenString string) (string, int, bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok == false {
			return nil, fmt.Errorf("Token signing method is not valid: %v", token.Header["alg"])
		}

		return []byte(c.Jwt.TokenSecret), nil
	})

	if err != nil {
		log.Println(err)
		return "", 0, false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		idFloat := claims["id"]
		id := int(idFloat.(float64))
		username := claims["username"]
		isAdmin := claims["is_admin"]
		return username.(string), id, isAdmin.(bool), nil
	}

	return "", 0, false, errors.New("Reading token claims failed")
}
