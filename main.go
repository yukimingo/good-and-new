package main

import (
	"good-and-new/controllers"
	"good-and-new/infra"
	"good-and-new/middlewares"
	"good-and-new/models"
	"good-and-new/repositories"
	"good-and-new/usecases"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()
	authRepository := repositories.NewAuthRepository(db)
	authUsecase := usecases.NewAuthUsecase(authRepository)
	authController := controllers.NewAuthController(authUsecase)

	NewsRepository := repositories.NewNewsRepository(db)
	NewsUsecase := usecases.NewNewsUsecase(NewsRepository)
	newsController := controllers.NewNewsController(NewsUsecase)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // 許可するオリジンを指定
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // クッキーや認証情報を許可
		MaxAge:           12 * time.Hour,
	}))
	authRouter := r.Group("/users")
	newsRouter := r.Group("/news")
	newsRouterWithAuth := r.Group("/news", middlewares.AuthMiddleware(authUsecase))
	{
		authRouter.POST("", authController.CreateUser)

		authRouter.POST("/login", authController.Login)

		authRouter.GET("", authController.FindAll)

		authRouter.GET("/email/:email", authController.FindUser)

		authRouter.GET("/:id", authController.FindUserById)

		authRouter.DELETE("/:id", authController.DeleteUser)
	}

	{
		newsRouter.GET("", newsController.FindAll)

		newsRouterWithAuth.GET("/:id", newsController.FindById)

		newsRouterWithAuth.POST("", newsController.Create)

		newsRouterWithAuth.DELETE("/:id", newsController.Delete)
	}

	authRouter.PUT("/:id", func(ctx *gin.Context) {
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

	r.Run("localhost:8080")
}
