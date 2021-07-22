package router

import (
	"fmt"
	"linker-fan/gal-anonim-server/server/handlers"

	"github.com/gin-gonic/gin"
)

func Run(port string, mode string) {
	router := setupRoutes()
	router.Run(fmt.Sprintf(":%s", port))
}

func setupRoutes() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	chat := r.Group("/chat")
	{
		chat.GET("/ws", handlers.ChatWebsocket)
	}

	return r
}
