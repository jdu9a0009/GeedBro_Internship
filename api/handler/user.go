package handler

import (
	"fmt"
	"net/http"
	"user/api/handler/response"
	"user/config"
	"user/models"
	"user/pkg/helper"
	"user/pkg/logger"

	"strconv"

	"github.com/gin-gonic/gin"
)

// SignInUser godoc
// @Router       /login [POST]
// @Summary      Login in User
// @Description  Login in User BASED ON GIVEN DATA
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data  body      models.LoginRequest  true  "user data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	err := c.ShouldBind(&req)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.storage.User().GetByLogin(c.Request.Context(), &models.LoginRequest{
		Username: req.Username,
	})
	if err != nil {
		fmt.Println("error User GetByLogin:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"not found login with:": req.Username})
		return
	}

	if req.Password != resp.Password {
		h.log.Error("error while binding:", logger.Error(err))
		res := response.ErrorResp{Code: "INVALID Password", Message: "invalid password"}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	m := make(map[string]interface{})
	m["username"] = resp.Username
	m["password"] = resp.Password
	m["phone"] = resp.Phone

	token, err := helper.GenerateJWT(m, config.TokenExpireTime, config.JWTSecretKey)
	c.JSON(http.StatusCreated, models.LoginRes{Token: token})
}

// SignUpUser godoc
// @Router       /signup [POST]
// @Summary      SignUp User
// @Description  SignUp User BASED ON GIVEN DATA
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateUser  true  "user data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) SignUp(c *gin.Context) {
	var user models.CreateUser
	err := c.ShouldBind(&user)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	if !helper.IsValidPhone(user.Phone) {
		c.JSON(http.StatusBadRequest, "invalid body,You should write phone numb uzbek version")
		return
	}

	_, err = h.storage.User().GetByLogin(c.Request.Context(), &models.LoginRequest{
		Username: user.Username,
	})

	if err != nil {
		resp, err := h.storage.User().CreateUser(c.Request.Context(), &user)
		if err != nil {
			h.log.Error("error Branch Create:", logger.Error(err))
			c.JSON(http.StatusInternalServerError, "internal server error")
			return
		}
		c.JSON(http.StatusCreated, response.CreateResponse{Message: "Succesfully created", Id: resp})
	}

	h.log.Error("error while binding:", logger.Error(err))
	res := response.ErrorResp{Code: "INVALID Username", Message: "This username  was used enter an other username"}
	c.JSON(http.StatusBadRequest, res)
	return
}

// CreateUser godoc
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
	err := c.ShouldBind(&user)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	// if !isvalid0hon(user.phone) {
	// c.JSON(http.StatusBadRequest, "invalid phone)
	// 	return
	// }

	resp, err := h.storage.User().CreateUser(c.Request.Context(), &user)
	if err != nil {
		h.log.Error("error User Create:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusCreated, response.CreateResponse{Message: "Succesfully created", Id: resp})
}

// Get user godoc
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
	id := c.Param("id")

	resp, err := h.storage.User().GetUser(c.Request.Context(), &models.IdRequest{Id: id})
	if err != nil {
		h.log.Error("error Branch Get:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetAllUsers godoc
// @Router       /user [GET]
// @Summary      GET  ALL Users
// @Description  get all branches based on limit, page and search by name
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

	resp, err := h.storage.User().GetAllUser(c.Request.Context(), &models.GetAllUserRequest{
		Page:  page,
		Limit: limit,
		Name:  c.Query("search"),
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
// @Router       /user/{id} [PUT]
// @Summary      UPDATE user BY ID
// @Description  UPDATES user BASED ON GIVEN DATA AND ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of user" format(uuid)
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

	user.ID = c.Param("id")
	resp, err := h.storage.User().UpdateUser(c.Request.Context(), &user)
	if err != nil {
		h.log.Error("error updating user:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully updated", "id": resp})
}

// DeleteUser godoc
// @Router       /user/{id} [DELETE]
// @Summary      DELETE user BY ID
// @Description  DELETES user BASED ON ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of user" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.storage.User().DeleteUser(c.Request.Context(), &models.IdRequest{Id: id})
	if err != nil {
		h.log.Error("error deleting user:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete branch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully deleted", "id": resp})
}
