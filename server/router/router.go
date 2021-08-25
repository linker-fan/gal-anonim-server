package router

import (
	"fmt"
	"linker-fan/gal-anonim-server/server/config"
	"linker-fan/gal-anonim-server/server/handlers"
	"linker-fan/gal-anonim-server/server/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func Run(port string, mode string, c *config.Config) {
	router, err := setupRoutes(c)
	if err != nil {
		log.Fatal(err)
	}
	router.Run(fmt.Sprintf(":%s", port))
}

func setupRoutes(c *config.Config) (*gin.Engine, error) {
	r := gin.Default()

	api, err := handlers.NewAPIWrapper(c)
	if err != nil {
		return nil, err
	}

	auth := r.Group("/users")
	{
		auth.POST("/register", api.Register)
		auth.POST("/login", api.Login)
	}

	protected := r.Group("/protected")
	protected.Use(middleware.JwtMiddleware())
	{
		protected.GET("/me", api.MeHandler)
		protected.POST("/refresh_token", api.RefreshTokenHandler)
		protected.POST("/pin", api.SetPinHandler)
	}

	room := r.Group("/room")
	room.Use(middleware.JwtMiddleware())
	{
		room.POST("", api.CreateRoomHandler)
		room.DELETE("/:uniqueRoomID", api.DeleteRoomHandler)
		room.PUT("/:uniqueRoomID", api.UpdateRoomDataHandler)
		room.GET("/:uniqueRoomID/members", api.GetRoomMembersHandler)
		room.POST("/:uniqueRoomID/add_member", api.AddMemberToTheRoomHandler)
		room.DELETE("/:uniqueRoomID/remove_member", api.RemoveMemberFromTheRoomHandler)
		room.DELETE("/:uniqueRoomID/leave", api.LeaveRoomHandler)
	}

	chat := r.Group("/chat")
	chat.Use(middleware.JwtMiddleware())
	chat.Use(middleware.ChatMiddleware())
	{
		chat.GET("/ws", api.ChatWebsocket)
		chat.POST("/:uniqueRoomID/upload", api.FileUploadHandler)
	}

	return r, nil
}
