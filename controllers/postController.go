package controllers

import (
	"blog-api/config"
	"blog-api/models"
	"blog-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePost(c *gin.Context) {
	var post models.Post
	post.Title = c.PostForm("title")
	post.Content = c.PostForm("content")

	// Lấy username từ middleware
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Kiểm tra vai trò user
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if user.Role != "VIP" && user.Role != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Xử lý file upload
	file, err := c.FormFile("image")
	if err == nil {
		fileName, err := utils.SaveUploadedFile(c, file, "public/uploads")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		post.ImageURL = "uploads/" + fileName
	}

	// Lấy user_id từ form
	userIDStr := c.PostForm("user_id")
	if userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}
		post.UserID = uint(userID)
	}

	// Lưu bài viết vào database
	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post created", "image": post.ImageURL})
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	config.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID, Username, Role") // Chỉ lấy các trường cần thiết
	}).Find(&posts)
	c.JSON(http.StatusOK, posts)
}
