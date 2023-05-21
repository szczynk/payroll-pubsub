package handler

import "github.com/gin-gonic/gin"

type EmployeeGinHandlerI interface {
	List(c *gin.Context)
	Add(c *gin.Context)
	Detail(c *gin.Context)
}
