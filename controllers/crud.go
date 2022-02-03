package controllers

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/malhotra-sidharth/url-shortener-go/models"
	"github.com/malhotra-sidharth/url-shortener-go/services"
)

type ICRUD interface {
	CreateShortUrl(ctx *gin.Context)
	RedirectToFullUrl(ctx *gin.Context)
	DeleteUrl(ctx *gin.Context)
}

type crud struct{}

func NewCrud() ICRUD {
	return &crud{}
}

// @ref: https://stackoverflow.com/a/51069900
func isUrlValid(input string) bool {
	validUrl, err := url.ParseRequestURI(input)
	if err != nil {
		return false
	}

	switch validUrl.Scheme {
	case "https":
	case "http":
	default:
		return false
	}

	_, err = net.LookupHost(validUrl.Host)
	if err != nil {
		return false
	}

	return true
}

func (crud *crud) CreateShortUrl(ctx *gin.Context) {
	var payload *models.CreateShortUrlPayload
	ctx.BindJSON(&payload)
	// validate url
	if !isUrlValid(payload.FullUrl) {
		ctx.JSON(http.StatusBadRequest, &gin.H{"message": "Invalid URL"})
		return
	}

	// insert entry
	id, err := services.Container.Shortener.Create(payload.FullUrl)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusConflict, &gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, &gin.H{"fullUrl": payload.FullUrl, "id": id})
}

func (crud *crud) RedirectToFullUrl(ctx *gin.Context) {
	id := ctx.Param("urlSlug")
	if id == "" {
		ctx.JSON(http.StatusNotFound, &gin.H{})
		return
	}

	// find full URL
	result, err := services.Container.Shortener.ResolveUrl(id)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusNotFound, &gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, &gin.H{
		"result": result,
	})
}

func (crud *crud) DeleteUrl(ctx *gin.Context) {
	id := ctx.Param("urlSlug")
	if id == "" {
		ctx.JSON(http.StatusNotFound, &gin.H{})
		return
	}

	// find full URL
	deletedCount, err := services.Container.Shortener.DeleteUrl(id)

	if err != nil || *deletedCount == 0 {
		fmt.Println(err)
		ctx.JSON(http.StatusNotFound, &gin.H{})
		return
	}

	ctx.JSON(http.StatusNoContent, &gin.H{})
}
