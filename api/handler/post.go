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

// CreatePost godoc
// @Router       /post [POST]
// @Summary      Creat new post
// @Description  creates a new post based on the given postname amd password
// @Tags         post
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreatePost  true  "post data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreatePost(c *gin.Context) {
	var post models.CreatePost
	fmt.Println("Before Handler", post)

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.storage.Post().CreatePost(c.Request.Context(), &post)
	if err != nil {
		h.log.Error("error Post Create:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusCreated, response.CreateResponse{Message: "Succesfully created", Id: resp})
}

// Get post godoc
// @Router       /post/{id} [GET]
// @Summary      GET BY ID
// @Description  get post by ID
// @Tags         post
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "post ID" format(uuid)
// @Success      200  {object}  models.Post
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetPost(c *gin.Context) {
	post := models.Post{}
	id := c.Param("id")
	ok, err := h.redisStorage.Cache().Get(c.Request.Context(), id, post)
	if err != nil {
		fmt.Println("not found staff in redis cache")
	}

	if ok {
		c.JSON(http.StatusOK, post)
		return
	}
	resp, err := h.storage.Post().GetPost(c.Request.Context(), &models.IdRequest{Id: id})
	if err != nil {
		h.log.Error("error Post Get:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": resp})

	err = h.redisStorage.Cache().Create(c.Request.Context(), id, resp, time.Hour)
	if err != nil {
		fmt.Println("error Post Create in redis cache:", err.Error())
	}
}

// GetAllPosts godoc
// @Router       /post [GET]
// @Summary      GET  ALL Posts
// @Description  get all posts based on limit, page and search by postname
// @Tags         post
// @Accept       json
// @Produce      json
// @Param   limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param   page         query     int        false  "page"          minimum(1)     default(1)
// @Param   search         query     string        false  "search"
// @Success      200  {object}  models.GetAllPost
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllPost(c *gin.Context) {
	h.log.Info("request GetAllPost")
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
	username := c.Query("search")
	resp, err := h.storage.Post().GetAllPost(c.Request.Context(), &models.GetAllPostRequest{
		Page:   page,
		Limit:  limit,
		Search: username,
	})

	if err != nil {
		h.log.Error("error  GetAllPost:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	h.log.Warn("response to GetAllPost")
	c.JSON(http.StatusOK, resp)
}

// UpdatePost godoc
// @Router       /post/{id} [PUT]
// @Summary      UPDATE post BY ID
// @Description  UPDATES post BASED ON GIVEN DATA AND ID
// @Tags         post
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of post" format(uuid)
// @Param        data  body      models.CreatePost  true  "post data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) UpdatePost(c *gin.Context) {
	var post models.Post
	err := c.ShouldBind(&post)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	post.ID = c.Param("id")
	resp, err := h.storage.Post().UpdatePost(c.Request.Context(), &post)
	if err != nil {
		h.log.Error("error updating post:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post successfully updated", "id": resp})
	err = h.redisStorage.Cache().Delete(c.Request.Context(), post.ID)
	if err != nil {
		fmt.Println("error Post Create in redis cache:", err.Error())
	}

}

// DeletePost godoc
// @Router       /post/{id} [DELETE]
// @Summary      DELETE post BY ID
// @Description  DELETES post BASED ON ID
// @Tags         post
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of post" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeletePost(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.Post().DeletePost(c.Request.Context(), &models.IdRequest{Id: id})
	if err != nil {
		h.log.Error("error deleting post:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete branch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post successfully deleted", "id": resp})
	err = h.redisStorage.Cache().Delete(c.Request.Context(), id)
	if err != nil {
		fmt.Println("error Post Create in redis cache:", err.Error())
	}
}
