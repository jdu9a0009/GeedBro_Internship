package handler

import (
	"fmt"
	"net/http"

	"user/models"
	"user/pkg/helper"
	"user/pkg/logger"

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
// @Param        data  body      models.CreateUserReq  true  "user data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreateUser(c *gin.Context) {
	var user models.CreateUserReq

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
