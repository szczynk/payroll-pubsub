package handler

import "github.com/gin-gonic/gin"

type APIKeyGinHandlerI interface {
	Generate(c *gin.Context)
}
