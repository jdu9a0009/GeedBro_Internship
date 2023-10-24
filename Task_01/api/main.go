package api

import (
	_ "user/api/docs"
	"user/api/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewServer(h *handler.Handler) *gin.Engine {
	r := gin.Default()
	//Users

	r.POST("/user", h.CreateUser)
	r.GET("/user/:id", h.GetUser)
	r.GET("/user", h.GetAllUser)
	r.PUT("/user/:id", h.UpdateUser)
	r.DELETE("/user/:id", h.DeleteUser)
	r.GET("/deleted_users", h.GetAllDeletedUser)

	//Posts
	r.POST("/post", h.CreatePost)
	r.GET("/post/:id", h.GetPost)
	r.GET("/post", h.GetAllPost)
	r.PUT("/post/:id", h.UpdatePost)
	r.DELETE("/post", h.DeletePost)
	r.GET("/my/post/:created_by", h.GetMyPost)
	r.GET("/deleted_posts", h.GetAllDeletedPost)

	// r.DELETE("/my/posts", h.getmypost)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return r
}
