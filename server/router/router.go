package router

import (
	"fmt"
	"linker-fan/gal-anonim-server/server/handlers"
	"linker-fan/gal-anonim-server/server/hub"
	"linker-fan/gal-anonim-server/server/middleware"

	"github.com/gin-gonic/gin"
)

func Run(port string, mode string) {
	//chat websocket
	wsServer := hub.NewHub()
	go wsServer.Run()
	router := setupRoutes(wsServer)
	router.Run(fmt.Sprintf(":%s", port))
}

func setupRoutes(hub *hub.Hub) *gin.Engine {
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
		protected.POST("/pin", handlers.SetPinHandler)
	}

	room := r.Group("/room")
	room.Use(middleware.JwtMiddleware())
	{
		room.POST("", handlers.CreateRoomHandler)
		room.DELETE("/:uniqueRoomID", handlers.DeleteRoomHandler)
		room.PUT("/:uniqueRoomID", handlers.UpdateRoomDataHandler)
		room.GET("/:uniqueRoomID/members", handlers.GetRoomMembersHandler)
		room.POST("/:uniqueRoomID/add_member", handlers.AddMemberToTheRoomHandler)
		room.DELETE("/:uniqueRoomID/remove_member", handlers.RemoveMemberFromTheRoomHandler)
		room.DELETE("/:uniqueRoomID/leave", handlers.LeaveRoomHandler)
	}

	chat := r.Group("/chat")
	chat.Use(middleware.JwtMiddleware())
	chat.Use(middleware.ChatMiddleware())
	{
		chat.GET("/:uniqueRoomID/ws", handlers.ChatWebsocket)
	}

	return r
}
