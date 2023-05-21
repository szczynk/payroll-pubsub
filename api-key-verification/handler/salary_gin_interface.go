package handler

import "github.com/gin-gonic/gin"

type SalaryGinHandlerI interface {
	List(c *gin.Context)
}
