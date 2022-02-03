package main

import (
	"github.com/gin-gonic/gin"
	"github.com/malhotra-sidharth/url-shortener-go/controllers"
	"github.com/malhotra-sidharth/url-shortener-go/services"
)

func registerRoutes(app *gin.Engine, container *services.ServiceContainer) {
	crud := controllers.NewCrud(container)
	app.GET("/", crud.HelloWorld)
	app.POST("/url", crud.CreateShortUrl)
	app.GET("/url/:urlSlug", crud.RedirectToFullUrl)
	app.DELETE("/url/:urlSlug", crud.DeleteUrl)
	app.GET("/url/:urlSlug/analytics", crud.GetAnalytics)
}
