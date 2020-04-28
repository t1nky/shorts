package controllers

import (
	"net/http"

	"shorts/database"
	h "shorts/helper"
	"shorts/models"

	"github.com/gin-gonic/gin"
)

// AddUser : Register a new user
func AddUser(c *gin.Context) {
	var userData models.AddUserData

	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, h.NewValidationError(userData, err))
		return
	}

	user := models.User{
		Name:     userData.Name,
		Password: userData.Password,
	}

	if dbc := database.DB.Create(&user); dbc.Error != nil {
		c.JSON(http.StatusBadRequest, h.NewResponseError(dbc.Error))
	} else {
		c.JSON(http.StatusCreated, h.NewResponseOK())
	}
}

// GetCurrentUser : Get currently authenticated user's information
func GetCurrentUser(c *gin.Context) {
	var user models.User

	userID := c.MustGet(gin.AuthUserKey).(uint64)

	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, h.NewResponseError(err))
	} else {
		// Remove password for *security reasons*
		// Remove shortlinks because we have /shorts
		c.JSON(http.StatusOK, h.NewResponseOkWithData(models.UserResponseData{
			ID:   user.ID,
			Name: user.Name,
		}))
	}
}
