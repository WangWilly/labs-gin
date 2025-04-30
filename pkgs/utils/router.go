package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDefaultRouter() *gin.Engine {
	router := gin.Default()

	// router.Use(gin.Recovery())
	// router.Use(gin.Logger())

	// https://stackoverflow.com/questions/32443738/setting-up-route-not-found-in-gin
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Endpoint not found"})
	})

	router.StaticFile("/favicon.ico", "./public/icon/favicon.ico")
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return router
}
