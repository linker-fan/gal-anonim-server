package router

import "github.com/gin-gonic/gin"

func Run(port string) {
	router := setupRoutes()
	router.Run(port)
}

func setupRoutes() *gin.Engine {
	r := gin.New()

	return r
}
