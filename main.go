package main

import (
	"good-and-new/controllers"
	"good-and-new/infra"
	"good-and-new/models"
	"good-and-new/repositories"
	"good-and-new/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()
	authRepository := repositories.NewAuthRepository(db)
	authUsecase := usecases.NewAuthUsecase(authRepository)
	authController := controllers.NewAuthController(authUsecase)

	n := &models.News{
		Title:       "title2",
		Description: "これはテスト2です",
		UserID:      7,
	}

	r := gin.Default()
	r.POST("/users", authController.CreateUser)

	r.GET("/users", authController.FindAll)

	r.GET("/users/:email", authController.FindUser)

	r.GET("/user/:id", authController.FindUserById)

	r.DELETE("/users/:id", authController.DeleteUser)

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
