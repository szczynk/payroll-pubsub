package handler

import (
	"fundamental-payroll-gin/helper/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingGinHandler struct{}

func NewPingGinHandler() PingGinHandlerI {
	return new(PingGinHandler)
}

func (h *PingGinHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, response.JSONRes{
		Status:  http.StatusOK,
		Message: "pong",
	})
}
