package router

import (
	"LYLChatBot/web/api"
	"LYLChatBot/web/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello KIT",
	})
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	r.Use(middleware.Auth()) //后续全部鉴权

	r.POST("send", api.Send)

	return r
}
