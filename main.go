package main

import (
	"good-and-new/infra"
	"good-and-new/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()
	// authRepository := repositories.NewAuthRepository(db)
	u := &models.User{
		Name:     "test3",
		Email:    "test3@example.com",
		Password: "password3",
	}
	n := &models.News{
		Title:       "title2",
		Description: "これはテスト2です",
		UserID:      7,
	}
	// updatedUser := &models.User{
	// 	Name:     "updated user",
	// 	Email:    "updated@example.com",
	// 	Password: "updatedPassword",
	// }

	r := gin.Default()
	r.POST("/users", func(ctx *gin.Context) {
		if err := db.Create(u).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		ctx.Status(http.StatusOK)
	})
	r.GET("/users", func(ctx *gin.Context) {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": &users})
	})
	r.GET("/users/:id", func(ctx *gin.Context) {
		var user models.User
		userId, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}
		if err := db.First(&user, "id = ?", userId).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": &user})
	})
	r.DELETE("/users/:id", func(ctx *gin.Context) {
		var user models.User
		userId, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}
		if err := db.First(&user, "id = ?", userId).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
			return
		}
		if err := db.Delete(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
			return
		}

		ctx.Status(http.StatusOK)
	})
	r.PUT("/users/:id", func(ctx *gin.Context) {
		var user models.User
		userId, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}
		if err := db.First(&user, "id = ?", userId).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
			return
		}
		user.Name = "updated user"
		user.Email = "update@example.com"
		user.Password = "updatedPassword"
		if err := db.Save(user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
			return
		}

		ctx.Status(http.StatusOK)
	})
	r.POST("/news", func(ctx *gin.Context) {
		if err := db.Create(n).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create news"})
			return
		}

		ctx.Status(http.StatusOK)
	})

	r.Run("localhost:8080")
}
