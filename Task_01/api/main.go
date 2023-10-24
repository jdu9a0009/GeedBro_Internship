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

	//Users

	r.POST("/user", helper.AuthMiddleware, h.CreateUser)
	r.GET("/user/:id", helper.AuthMiddleware, h.GetUser)
	r.GET("/user", helper.AuthMiddleware, h.GetAllUser)
	r.PUT("/user", helper.AuthMiddleware, h.UpdateUser)
	r.DELETE("/user", helper.AuthMiddleware, h.DeleteUser)
	r.GET("/deleted_users", helper.AuthMiddleware, h.GetAllDeletedUser)

	//Posts
	r.POST("/post", helper.AuthMiddleware, h.CreatePost)
	r.GET("/post/:id", helper.AuthMiddleware, h.GetPost)
	r.GET("/post", helper.AuthMiddleware, h.GetAllPost)
	r.PUT("/post", helper.AuthMiddleware, h.UpdatePost)
	r.DELETE("/post", helper.AuthMiddleware, h.DeletePost)
	r.GET("/my/post/:created_by", helper.AuthMiddleware, h.GetMyPost)
	r.GET("/deleted_posts", helper.AuthMiddleware, h.GetAllDeletedPost)

	// r.DELETE("/my/posts", h.getmypost)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return r
}
