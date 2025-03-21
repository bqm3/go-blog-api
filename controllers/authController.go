package controllers

import (
	"blog-api/config"
	"blog-api/models"
	"blog-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)

	// Mã hóa mật khẩu trước khi lưu
	hashedPassword, err := utils.HashPassword((user.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = hashedPassword
	config.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

func LoginUser(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	var dbUser models.User

	config.DB.Where("username = ?", user.Username).First(&dbUser)
	if dbUser.ID == 0 || !utils.CheckPassword(dbUser.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token, _ := utils.GenerateToken(dbUser.Username)
	c.JSON(http.StatusOK, gin.H{"token": token, "data": dbUser})
}
