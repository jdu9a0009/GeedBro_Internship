package handler

import (
	"fmt"
	"net/http"
	"user/api/handler/response"
	"user/models"
	"user/pkg/logger"

	"time"

	"strconv"

	"github.com/gin-gonic/gin"
)

// CreatePostComment godoc
// @Security ApiKeyAuth
// @Router       /post_comment [POST]
// @Summary      Creat new post_comment
// @Description  creates a new post_comment based on the given post_commentname ad password
// @Tags         post_comment
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreatePostComment  true  "post_comment data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreatePostComment(c *gin.Context) {
	var post_comment models.CreatePostComment
	fmt.Println("Before Handler", post_comment)

	if err := c.ShouldBindJSON(&post_comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.storage.PostComment().CreatePostComment(c, &post_comment)
	if err != nil {
		h.log.Error("error PostComment Create:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusCreated, response.CreateResponse{Message: "Succesfully created", Id: resp})
}

// Get post_comment godoc
// @Security ApiKeyAuth
// @Router       /post_comment/{id} [GET]
// @Summary      GET BY ID
// @Description  get post_comment by PostCommentID
// @Tags         post_comment
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "post_comment ID" format(uuid)
// @Success      200  {object}  models.PostComment
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetPostComment(c *gin.Context) {
	post_comment := models.PostComment{}
	id := c.Param("id")
	ok, err := h.redisStorage.Cache().Get(c.Request.Context(), id, post_comment)
	if err != nil {
		fmt.Println("not found staff in redis cache")
	}

	if ok {
		c.JSON(http.StatusOK, post_comment)
		return
	}
	resp, err := h.storage.PostComment().GetPostComment(c, &models.IdRequest{Id: id})
	if err != nil {
		h.log.Error("error PostComment Get:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": resp})

	err = h.redisStorage.Cache().Create(c.Request.Context(), id, resp, time.Hour)
	if err != nil {
		fmt.Println("error PostComment Create in redis cache:", err.Error())
	}
}

// GetAllPostComments godoc
// @Security ApiKeyAuth
// @Router       /post_comment [GET]
// @Summary      GET  ALL PostComments
// @Description  get all post_comments based on limit, page and search by post_commentname
// @Tags         post_comment
// @Accept       json
// @Produce      json
// @Param   limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param   page         query     int        false  "page"          minimum(1)     default(1)
// @Param   search         query     string        false  "search"
// @Success      200  {object}  models.GetAllPostComment
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllPostComment(c *gin.Context) {
	h.log.Info("request GetAllPostComment")
	page, err := strconv.Atoi(c.DefaultQuery("page", "fmt.sprintf(`%d`,cfg.DefaultPage)"))
	if err != nil {
		h.log.Error("error getting page:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid page param")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		h.log.Error("error getting limit:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid limit param")
		return
	}
	Post_id := c.Query("search")
	resp, err := h.storage.PostComment().GetAllPostComment(c, &models.GetAllPostCommentRequest{
		Page:    page,
		Limit:   limit,
		Post_id: Post_id,
	})

	if err != nil {
		h.log.Error("error  GetAllPostComment:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	h.log.Warn("response to GetAllPostComment")
	c.JSON(http.StatusOK, resp)
}

// UpdatePostComment godoc
// @Security ApiKeyAuth
// @Router       /post_comment [PUT]
// @Summary      UPDATE post_comment BY ID
// @Description  UPDATES post_comment BASED ON GIVEN DATA AND ID
// @Tags         post_comment
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of post_comment" format(uuid)
// @Param        data  body      models.UpdatePostComment true  "post_comment data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) UpdatePostComment(c *gin.Context) {
	var post_comment models.UpdatePostComment
	err := c.ShouldBind(&post_comment)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.storage.PostComment().UpdatePostComment(c, &post_comment)
	if err != nil {
		h.log.Error("error updating post_comment:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post_comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PostComment successfully updated", "id": resp})
	err = h.redisStorage.Cache().Delete(c.Request.Context(), post_comment.ID)
	if err != nil {
		fmt.Println("error PostComment Create in redis cache:", err.Error())
	}

}

// DeletePostComment godoc
// @Security ApiKeyAuth
// @Router       /post_comment [DELETE]
// @Summary      DELETE post_comment BY ID
// @Description  DELETES post_comment BASED ON ID
// @Tags         post_comment
// @Accept       json
// @Produce      json
// @Param        data  body      models.DeletePostCommentRequest  true  "post_comment data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeletePostComment(c *gin.Context) {
	var post_comment models.DeletePostCommentRequest
	err := c.ShouldBindJSON(&post_comment)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	resp, err := h.storage.PostComment().DeletePostComment(c, &models.DeletePostCommentRequest{
		Id: post_comment.Id,
	})
	if err != nil {
		h.log.Error("error deleting post_comment:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete branch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PostComment successfully deleted", "id": resp})
	err = h.redisStorage.Cache().Delete(c.Request.Context(), post_comment.Id)
	if err != nil {
		fmt.Println("error PostComment Create in redis cache:", err.Error())
	}
}

// Get post_comment godoc
// @Security ApiKeyAuth
// @Router       /my/post_comment/{created_by} [GET]
// @Summary      GET BY ID
// @Description  get post_comment by ID
// @Tags         post_comment
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.GetAllPostComment
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetMyPostComment(c *gin.Context) {

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		h.log.Error("error getting page:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid page param")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		h.log.Error("error getting limit:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid limit param")
		return
	}
	post_id := c.Query("search")
	resp, err := h.storage.PostComment().GetMyPostComment(c, &models.GetAllPostCommentRequest{
		Page:    page,
		Limit:   limit,
		Post_id: post_id,
	})
	if err != nil {
		h.log.Error("error  GetAllPostComment:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	h.log.Warn("response to GetAllPostComment")
	c.JSON(http.StatusOK, resp)
}

// GetAllPostComments godoc
// @Security ApiKeyAuth
// @Router       /deleted_post_comments [GET]
// @Summary      GET  ALL PostComments
// @Description  get all post_comments based on limit, page and search by post_commentname
// @Tags         post_comment
// @Accept       json
// @Produce      json
// @Param   limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param   page         query     int        false  "page"          minimum(1)     default(1)
// @Param   search         query     string        false  "search"
// @Success      200  {object}  models.GetAllPostComment
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllDeletedPostComment(c *gin.Context) {
	h.log.Info("request GetAllPostComment")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		h.log.Error("error getting page:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid page param")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		h.log.Error("error getting limit:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid limit param")
		return
	}
	post_id := c.Query("search")
	resp, err := h.storage.PostComment().GetAllDeletedPostComment(c, &models.GetAllPostCommentRequest{
		Page:    page,
		Limit:   limit,
		Post_id: post_id,
	})

	if err != nil {
		h.log.Error("error  GetAllPostComment:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	h.log.Warn("response to GetAllPostComment")
	c.JSON(http.StatusOK, resp)
}
