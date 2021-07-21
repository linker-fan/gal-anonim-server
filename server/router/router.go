package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Run(port string, mode string) {
	router := setupRoutes()
	router.Run(fmt.Sprintf(":%s", port))
}

func setupRoutes() *gin.Engine {
	r := gin.New()

	return r
}
