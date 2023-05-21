package server

import (
	"fundamental-payroll-gin/handler"
	"fundamental-payroll-gin/helper/response"
	"fundamental-payroll-gin/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewRouter(
	debug bool,
	logger *zerolog.Logger,
	employeeHandler handler.EmployeeGinHandlerI,
	payrollHandler handler.PayrollGinHandlerI,
	salaryHandler handler.SalaryGinHandlerI,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	if debug {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.JSONRes{
			Status:  http.StatusOK,
			Message: "pong",
		})
	})

	employeeRouter := router.Group("/employees")
	{
		employeeRouter.GET("", employeeHandler.List)
		employeeRouter.POST("", employeeHandler.Add)
		employeeRouter.GET("/:id", employeeHandler.Detail)
	}

	payrollRouter := router.Group("/payrolls")
	{
		payrollRouter.GET("", payrollHandler.List)
		payrollRouter.POST("", payrollHandler.Add)
		payrollRouter.GET("/:id", payrollHandler.Detail)
	}

	salaryRouter := router.Group("/salaries")
	{
		salaryRouter.GET("", salaryHandler.List)
	}

	return router
}
