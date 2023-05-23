package handler

import "github.com/gin-gonic/gin"

type PingGinHandlerI interface {
	Ping(c *gin.Context)
}
