package middleware

import (
	"linker-fan/gal-anonim-server/server/auth"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//JwtMiddleware - not sure wether to use cookies or return tokenstrings
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		cookie, err := c.Request.Cookie("jwtToken")
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		tokenString := cookie.Value
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
}
