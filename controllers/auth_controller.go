package controllers

import (
	"good-and-new/dto"
	"good-and-new/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IAuthController interface {
	FindAll(ctx *gin.Context)
	FindUser(ctx *gin.Context)
	FindUserById(ctx *gin.Context)
	CreateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type AuthController struct {
	usecase usecases.IAuthUsecase
}

func NewAuthController(u usecases.IAuthUsecase) IAuthController {
	return &AuthController{usecase: u}
}

func (c *AuthController) FindAll(ctx *gin.Context) {
	users, err := c.usecase.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": users})
}

func (c *AuthController) FindUser(ctx *gin.Context) {
	email := ctx.Param("email")
	user, err := c.usecase.FindUser(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

func (c *AuthController) CreateUser(ctx *gin.Context) {
	var userInput dto.SignupInput
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.usecase.CreateUser(userInput); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *AuthController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	err = c.usecase.DeleteUser(id)
	if err != nil {
		if err.Error() == "User not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *AuthController) FindUserById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	user, err := c.usecase.FindUserById(id)
	if err != nil {
		if err.Error() == "User not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}
