package handler

import "github.com/gin-gonic/gin"

type PayrollGinHandlerI interface {
	List(c *gin.Context)
	Add(c *gin.Context)
	Detail(c *gin.Context)
}
