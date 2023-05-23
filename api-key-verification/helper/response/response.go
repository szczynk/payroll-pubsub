package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONRes struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewJSONRes(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, JSONRes{
		Status:  statusCode,
		Message: http.StatusText(statusCode),
		Data:    data,
	})
}

func NewJSONResErr(c *gin.Context, statusCode int, err string) {
	c.JSON(statusCode, JSONRes{
		Status:  statusCode,
		Message: http.StatusText(statusCode),
		Error:   err,
	})
}

// func NewJSONResMsg(c *gin.Context, statusCode int, msg string, data any) {
// 	c.JSON(statusCode, JSONRes{
// 		Status:  statusCode,
// 		Message: msg,
// 		Data:    data,
// 	})
// }

// func NewJSONResMsgErr(c *gin.Context, statusCode int, msg string, err string) {
// 	c.JSON(statusCode, JSONRes{
// 		Status:  statusCode,
// 		Message: msg,
// 		Error:   err,
// 	})
// }
