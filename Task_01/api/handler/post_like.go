package handler

import (
	"fmt"
	"net/http"
	"strings"
	"user/api/handler/response"
	"user/models"
	"user/pkg/logger"

	"time"

	"github.com/gin-gonic/gin"
)

// CreatePostLike godoc
// @Security ApiKeyAuth
// @Router       /post_like [POST]
// @Summary      Creat new post
// @Description  creates a new post based on the given postname and password
// @Tags         post_like
// @Accept       json
// @Produce      json
// @Param        data  body      models.PostLikes  true  "post data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreatePostLike(c *gin.Context) {
	var post models.PostLikes
	fmt.Println("Before Handler", post)

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.storage.PostLike().CreatePostLike(c, &post)
	if err != nil {
		h.log.Error("error PostLike Create:", logger.Error(err))
		if strings.Contains(err.Error(), "like already exists") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "like already exists for the given post and user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		fmt.Println("Error:", err) // Print the error message to the console
		return
	}

	c.JSON(http.StatusCreated, response.CreateResponse{Message: "Successfully created", Id: resp})
}

// Get post godoc
// @Security ApiKeyAuth
// @Router       /post_like/{id} [GET]
// @Summary      GET BY ID
// @Description  get post by PostLikeID
// @Tags         post_like
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Like ID" format(uuid)
// @Success      200  {object}  models.PostLikes
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetPostLike(c *gin.Context) {
	post := models.PostLikes{}
	post_id := c.Param("id")
	ok, err := h.redisStorage.Cache().Get(c.Request.Context(), post_id, post)
	if err != nil {
		fmt.Println("not found staff in redis cache")
	}

	if ok {
		c.JSON(http.StatusOK, post)
		return
	}
	resp, err := h.storage.PostLike().GetPostLikes(c, &models.PostLikes{Post_Id: post_id})
	if err != nil {
		h.log.Error("error PostLike Get:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": resp})

	err = h.redisStorage.Cache().Create(c.Request.Context(), post_id, resp, time.Hour)
	if err != nil {
		fmt.Println("error PostLike Create in redis cache:", err.Error())
	}
}

// DeletePostLike godoc
// @Security ApiKeyAuth
// @Router       /post_like [DELETE]
// @Summary      DELETE post BY ID
// @Description  DELETES post BASED ON ID
// @Tags         post_like
// @Accept       json
// @Produce      json
// @Param        data  body      models.DeletePostLikeRequest  true  "post data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeletePostLike(c *gin.Context) {
	var post models.DeletePostLikeRequest
	err := c.ShouldBindJSON(&post)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	resp, err := h.storage.PostLike().DeletePostLike(c, &models.DeletePostLikeRequest{
		Post_Id: post.Post_Id,
	})
	if err != nil {
		h.log.Error("error deleting post:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Like"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PostLike successfully deleted", "id": resp})
	err = h.redisStorage.Cache().Delete(c.Request.Context(), post.Post_Id)
	if err != nil {
		fmt.Println("error PostLike Create in redis cache:", err.Error())
	}
}
