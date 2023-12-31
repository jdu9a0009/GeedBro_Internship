package handler

import (
	"fmt"
	"net/http"
	"time"
	"user/api/handler/response"
	"user/models"
	"user/pkg/helper"
	"user/pkg/logger"

	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Security ApiKeyAuth
// @Router       /user [POST]
// @Summary      CREATES User
// @Description  CREATES User BASED ON GIVEN DATA
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateUser  true  "user data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreateUser(c *gin.Context) {
	var user models.CreateUser
	fmt.Println("Before Handler", user)

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("After Handler", user)
	hashedPass, err := helper.GeneratePasswordHash(user.Password)
	if err != nil {
		h.log.Error("error while generating hash password:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}
	user.Password = string(hashedPass)

	resp, err := h.storage.User().CreateUser(c.Request.Context(), &user)
	if err != nil {
		h.log.Error("error User Create:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusCreated, response.CreateResponse{Message: "Succesfully created", Id: resp})
}

// Get user godoc
// @Security ApiKeyAuth
// @Router       /user/{id} [GET]
// @Summary      GET BY ID
// @Description  get user by ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "user ID" format(uuid)
// @Success      200  {object}  models.User
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetUser(c *gin.Context) {
	user := models.User{}
	id := c.Param("id")
	ok, err := h.redisStorage.Cache().Get(c.Request.Context(), id, user)
	if err != nil {
		fmt.Println("not found staff in redis cache")
	}

	if ok {
		c.JSON(http.StatusOK, user)
		return
	}
	resp, err := h.storage.User().GetUser(c.Request.Context(), &models.IdRequest{Id: id})
	if err != nil {
		h.log.Error("error User Get:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": resp})

	err = h.redisStorage.Cache().Create(c.Request.Context(), id, resp, time.Hour)
	if err != nil {
		fmt.Println("error User Create in redis cache:", err.Error())
	}
}

// GetAllUsers godoc
// @Security ApiKeyAuth
// @Router       /user [GET]
// @Summary      GET  ALL Users
// @Description  get all users based on limit, page and search by username
// @Tags         user
// @Accept       json
// @Produce      json
// @Param   limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param   page         query     int        false  "page"          minimum(1)     default(1)
// @Param   search         query     string        false  "search"
// @Success      200  {object}  models.GetAllUser
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllUser(c *gin.Context) {
	h.log.Info("request GetAllUser")
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

	resp, err := h.storage.User().GetAllUser(c.Request.Context(), &models.GetAllUserRequest{
		Page:     page,
		Limit:    limit,
		UserName: c.Query("search"),
	})
	if err != nil {
		h.log.Error("error Branch GetAllBUser:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	h.log.Warn("response to GetAllUser")
	c.JSON(http.StatusOK, resp)
}

// GetAllDeletedUsers godoc
// @Security ApiKeyAuth
// @Router       /deleted_users [GET]
// @Summary      GET  ALL Users
// @Description  get all users based on limit, page and search by name
// @Tags         user
// @Accept       json
// @Produce      json
// @Param   limit         query     int        false  "limit"          minimum(1)     default(10)
// @Param   page         query     int        false  "page"          minimum(1)     default(1)
// @Param   search         query     string        false  "search"
// @Success      200  {object}  models.GetAllUser
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllDeletedUser(c *gin.Context) {
	h.log.Info("request GetAllUser")
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

	resp, err := h.storage.User().GetAllDeletedUser(c.Request.Context(), &models.GetAllUserRequest{
		Page:     page,
		Limit:    limit,
		UserName: c.Query("search"),
	})
	if err != nil {
		h.log.Error("error Branch GetAllBUser:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	h.log.Warn("response to GetAllUser")
	c.JSON(http.StatusOK, resp)
}

// UpdateUser godoc
// @Security ApiKeyAuth
// @Router       /user [PUT]
// @Summary      UPDATE user BY ID
// @Description  UPDATES user BASED ON GIVEN DATA AND ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateUser  true  "user data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) UpdateUser(c *gin.Context) {
	var user models.User
	err := c.ShouldBind(&user)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.storage.User().UpdateUser(c, &user)
	if err != nil {
		h.log.Error("error updating user:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully updated", "id": resp})
	err = h.redisStorage.Cache().Delete(c.Request.Context(), user.ID)
	if err != nil {
		fmt.Println("error User Create in redis cache:", err.Error())
	}

}

// DeleteUser godoc
// @Security ApiKeyAuth
// @Router       /user [DELETE]
// @Summary      DELETE user BY ID
// @Description  DELETES user BASED ON ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.User().DeleteUser(c, &models.IdRequest{Id: id})
	if err != nil {
		h.log.Error("error deleting user:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete branch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully deleted", "id": resp})
	err = h.redisStorage.Cache().Delete(c.Request.Context(), id)
	if err != nil {
		fmt.Println("error User Create in redis cache:", err.Error())
	}
}
