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

	//Postlikes
	r.POST("/post_like", helper.AuthMiddleware, h.CreatePostLike)
	r.GET("/post_like/:id", helper.AuthMiddleware, h.GetPostLike)
	r.DELETE("/post_like", helper.AuthMiddleware, h.DeletePostLike)

	//Commentlikes
	r.POST("/comment_like", helper.AuthMiddleware, h.CreateCommentLike)
	r.GET("/comment_like/:id", helper.AuthMiddleware, h.GetCommentLike)
	r.DELETE("/comment_like", helper.AuthMiddleware, h.DeleteCommentLike)

	//Post_Comments
	r.POST("/post_comment", helper.AuthMiddleware, h.CreatePostComment)
	r.GET("/post_comment/:id", helper.AuthMiddleware, h.GetPostComment)
	r.GET("/post_comment", helper.AuthMiddleware, h.GetAllPostComment)
	r.PUT("/post_comment", helper.AuthMiddleware, h.UpdatePostComment)
	r.DELETE("/post_comment", helper.AuthMiddleware, h.DeletePostComment)
	r.GET("/my/post_comment/:created_by", helper.AuthMiddleware, h.GetMyPostComment)
	r.GET("/deleted_post_comments", helper.AuthMiddleware, h.GetAllDeletedPostComment)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return r
}
