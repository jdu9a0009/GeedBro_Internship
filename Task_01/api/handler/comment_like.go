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

// CreateCommentLike godoc
// @Security ApiKeyAuth
// @Router       /comment_like [POST]
// @Summary      Creat new comment
// @Description  creates a new comment based on the given commentname and password
// @Tags         comment_like
// @Accept       json
// @Produce      json
// @Param        data  body      models.CommentLikes  true  "comment data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreateCommentLike(c *gin.Context) {
	var comment models.CommentLikes
	fmt.Println("Before Handler", comment)

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.storage.CommentLike().CreateCommentLike(c, &comment)
	if err != nil {
		h.log.Error("error CommentLike Create:", logger.Error(err))
		if strings.Contains(err.Error(), "like already exists") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "like already exists for the given comment and user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		fmt.Println("Error:", err) // Print the error message to the console
		return
	}

	c.JSON(http.StatusCreated, response.CreateResponse{Message: "Successfully created", Id: resp})
}

// Get comment godoc
// @Security ApiKeyAuth
// @Router       /comment_like/{id} [GET]
// @Summary      GET BY ID
// @Description  get comment by CommentLikeID
// @Tags         comment_like
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Like ID" format(uuid)
// @Success      200  {object}  models.CommentLikes
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetCommentLike(c *gin.Context) {
	comment := models.CommentLikes{}
	comment_id := c.Param("id")
	ok, err := h.redisStorage.Cache().Get(c.Request.Context(), comment_id, comment)
	if err != nil {
		fmt.Println("not found staff in redis cache")
	}

	if ok {
		c.JSON(http.StatusOK, comment)
		return
	}
	resp, err := h.storage.CommentLike().GetCommentLikes(c, &models.CommentLikes{Comment_Id: comment_id})
	if err != nil {
		h.log.Error("error CommentLike Get:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": resp})

	err = h.redisStorage.Cache().Create(c.Request.Context(), comment_id, resp, time.Hour)
	if err != nil {
		fmt.Println("error CommentLike Create in redis cache:", err.Error())
	}
}

// DeleteCommentLike godoc
// @Security ApiKeyAuth
// @Router       /comment_like [DELETE]
// @Summary      DELETE comment BY ID
// @Description  DELETES comment BASED ON ID
// @Tags         comment_like
// @Accept       json
// @Produce      json
// @Param        data  body      models.DeleteCommentLikeRequest  true  "comment data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteCommentLike(c *gin.Context) {
	var comment models.DeleteCommentLikeRequest
	err := c.ShouldBindJSON(&comment)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	resp, err := h.storage.CommentLike().DeleteCommentLike(c, &models.DeleteCommentLikeRequest{
		Comment_Id: comment.Comment_Id,
	})
	if err != nil {
		h.log.Error("error deleting comment:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Like"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "CommentLike successfully deleted", "id": resp})
	err = h.redisStorage.Cache().Delete(c.Request.Context(), comment.Comment_Id)
	if err != nil {
		fmt.Println("error CommentLike Create in redis cache:", err.Error())
	}
}
