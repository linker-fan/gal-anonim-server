package middleware

import (
	"errors"
	"linker-fan/gal-anonim-server/server/auth"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//JwtMiddleware - not sure wether to use cookies or return tokenstrings
func JwtMiddleware() gin.HandlerFunc {
	return authmiddleware
}

func authmiddleware(c *gin.Context) {
	var tokenString string

	cookie, err := c.Request.Cookie("jwtToken")
	if err == http.ErrNoCookie {
		//if the cookie is not set, check the header
		tokenString, err = getTokenFromHeader(c.Request)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
	} else {
		tokenString = cookie.Value
	}

	username, id, isAdmin, err := auth.IsTokenValid(tokenString)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Not Authorized",
		})
		c.Abort()
		return
	} else {
		c.Set("id", id)
		c.Set("username", username)
		c.Set("is_admin", isAdmin)
		c.Next()
	}
}

func getTokenFromHeader(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) != 2 {
		return "", errors.New("Not Valid")
	}
	tokenString := strings.TrimSpace(splitToken[1])
	return tokenString, nil
}
