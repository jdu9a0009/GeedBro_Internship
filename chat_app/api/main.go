package api

import (
	_ "user/api/docs"
	"user/api/handler"
	"user/pkg/helper"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewServer(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	//Auth sign up and login
	r.POST("/auth/login", h.Login)
	r.POST("/auth/sign-up", h.SignUp)

	r.POST("/ws/createRoom", h.CreateRoom)
	r.GET("/ws/joinRoom/:roomId", h.JoinRoom)
	r.GET("/ws/getRooms", h.GetRooms)
	r.GET("/ws/getClients/:roomId", h.GetClients)

	r.POST("/user", helper.AuthMiddleware, h.CreateUser)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return r
}
