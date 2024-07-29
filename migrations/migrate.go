package main

import (
	"good-and-new/infra"
	"good-and-new/models"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()
	if err := db.AutoMigrate(&models.User{}, &models.News{}); err != nil {
		panic("Failed to migrate database")
	}
}
