package handler

import (
	"fmt"
	"net/http"
	"user/config"
	"user/models"
	"user/pkg/helper"
	"user/pkg/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// SignUp User godoc
// @Router       /auth/sign-up [POST]
// @Summary      Sign Up User
// @Description  Sign Up User BASED ON GIVEN DATA
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
	err := c.ShouldBindJSON(&user)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	hashedPass, err := helper.GeneratePasswordHash(user.Password)
	if err != nil {
		h.log.Error("error while generating hash password:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}
	user.Password = string(hashedPass)

	resp, err := h.storage.User().CreateUser(c.Request.Context(), &user)
	if err != nil {
		fmt.Println("error User Create:", err.Error())
		c.JSON(http.StatusInternalServerError, "username is already used, enter another one")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created", "id": resp})
}

// Login User godoc
// @Router       /auth/login [POST]
// @Summary      login User
// @Description  login User BASED ON GIVEN DATA
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
	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fields in body"})
		return
	}

	resp, err := h.storage.User().GetByLogin(c, &models.LoginRequest{
		Username: req.Username,
	})
	if err != nil {
		h.log.Error("error get by username:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "not found username"})
		return
	}

	// Compare hashed password with plain text password
	err = helper.ComparePasswords([]byte(resp.Password), []byte(req.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "login or password didn't match"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "password comparison failed"})
		}
		return
	}

	m := make(map[string]interface{})
	m["user_id"] = resp.UserId

	token, _ := helper.GenerateJWT(m, config.TokenExpireTime, config.JWTSecretKey)
	c.JSON(http.StatusOK, models.LoginRespond{Token: token})
}
