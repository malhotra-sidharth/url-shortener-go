package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/malhotra-sidharth/url-shortener-go/services"
)

func main() {
	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Content-Type"},
		AllowMethods:    []string{"GET", "POST", "DELETE"},
	}))
	container := services.RegisterServiceContainer()
	defer services.DbDisconnect()
	registerRoutes(app, container)
	app.Run()
}
