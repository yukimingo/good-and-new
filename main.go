package main

import (
	"good-and-new/infra"
	"good-and-new/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()
	// authRepository := repositories.NewAuthRepository(db)
	u := &models.User{
		Name:     "test",
		Email:    "test@example.com",
		Password: "password",
	}

	r := gin.Default()
	r.POST("/users", func(ctx *gin.Context) {
		if err := db.Create(u).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		ctx.Status(http.StatusOK)
	})

	r.Run("localhost:8080")
}
