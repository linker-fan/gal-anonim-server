package middleware

import (
	"github.com/gin-gonic/gin"
)

//JwtMiddleware - not sure wether to use cookies or return tokenstrings
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*
			cookie, err := c.Request.Cookie("jwtToken")
			if err != nil {
				c.Status(http.StatusUnauthorized)
				c.Abort()
				return
			}

			tokenString := cookie.Value
		*/
	}
}
