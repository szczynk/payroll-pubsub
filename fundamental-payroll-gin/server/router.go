package server

import (
	"fundamental-payroll-gin/handler"
	"fundamental-payroll-gin/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewRouter(
	debug bool,
	apiVerificationURL string,
	logger *zerolog.Logger,
	pingHandler handler.PingGinHandlerI,
	employeeHandler handler.EmployeeGinHandlerI,
	payrollHandler handler.PayrollGinHandlerI,
	salaryHandler handler.SalaryGinHandlerI,
	apiKeyVerificationHandler handler.APIKeyGinHandlerI,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	if debug {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())

	router.GET("/ping", pingHandler.Ping)

	router.POST("/generate", apiKeyVerificationHandler.Generate)

	apiKeyMiddleware := middleware.APIKey(apiVerificationURL)

	employeeRouter := router.Group("/employees")
	employeeRouter.Use(apiKeyMiddleware)
	{
		employeeRouter.GET("", employeeHandler.List)
		employeeRouter.POST("", employeeHandler.Add)
		employeeRouter.GET("/:id", employeeHandler.Detail)
	}

	payrollRouter := router.Group("/payrolls")
	payrollRouter.Use(apiKeyMiddleware)
	{
		payrollRouter.GET("", payrollHandler.List)
		payrollRouter.POST("", payrollHandler.Add)
		payrollRouter.GET("/:id", payrollHandler.Detail)
	}

	salaryRouter := router.Group("/salaries")
	salaryRouter.Use(apiKeyMiddleware)
	{
		salaryRouter.GET("", salaryHandler.List)
	}

	return router
}
