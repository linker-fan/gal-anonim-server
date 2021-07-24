package router

import (
	"fmt"
	"linker-fan/gal-anonim-server/server/handlers"
	"linker-fan/gal-anonim-server/server/middleware"

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

	auth := r.Group("/users")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	protected := r.Group("/protected")
	protected.Use(middleware.JwtMiddleware())
	{
		protected.GET("/me", handlers.MeHandler)
		protected.POST("/refresh_token", handlers.RefreshTokenHandler)
	}

	room := r.Group("/room")
	//room.Use(middleware.JwtMiddleware())
	{
		room.POST("", handlers.CreateRoomHandler)
		room.DELETE("", handlers.DeleteRoomHandler)
	}

	chat := r.Group("/chat")
	{
		chat.GET("/ws", handlers.ChatWebsocket)
	}

	return r
}
