package main

import (
	"github.com/gin-gonic/gin"
	"github.com/malhotra-sidharth/url-shortener-go/controllers"
)

func registerRoutes(app *gin.Engine) {
	crud := controllers.NewCrud()
	app.POST("/url", crud.CreateShortUrl)
	app.GET("/url/:urlSlug", crud.RedirectToFullUrl)
	app.DELETE("/url/:urlSlug", crud.DeleteUrl)
	app.GET("/url/:urlSlug/analytics", crud.GetAnalytics)
}
