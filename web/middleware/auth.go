package middleware

import (
	"LYLChatBot/conf"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
)

// Auth middleware
func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.Request.Header.Get("Authorization")

		ip := net.ParseIP(context.ClientIP())
		for _, auth := range conf.ConfigureInstance.APIAuth {
			if auth.Token == authHeader {
				_, cidr, err := net.ParseCIDR(auth.Source)
				if err == nil {
					if cidr.Contains(ip) {
						context.Next()
						return
					} else {
						continue
					}
				} else {
					log.Error(err)
				}

			}
		}
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		context.Abort()
	}
}
