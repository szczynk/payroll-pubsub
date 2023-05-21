package main

import (
	"api-key-verification/config"
	"api-key-verification/handler"
	"api-key-verification/helper/logger"
	"api-key-verification/helper/response"
	"api-key-verification/middleware"
	"api-key-verification/server"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config, confErr := config.LoadConfig()
	if confErr != nil {
		log.Fatalf("load config err:%s", confErr)
	}

	logger := logger.New(config.Debug)

	apiKeyHandler := handler.NewAPIKeyGinHandler(config.Passphrase)
	employeeHandler := handler.NewEmployeeGinHandler(config.APIURL)
	payrollHandler := handler.NewPayrollGinHandler(config.APIURL)
	salaryHandler := handler.NewSalaryGinHandler(config.APIURL)

	router := gin.New()

	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())
	// router.Use(middleware.APIKey(config.Passphrase))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.JSONRes{
			Status:  http.StatusOK,
			Message: "pong",
		})
	})

	router.POST("/generate", apiKeyHandler.Generate)

	employeeRouter := router.Group("/employees")
	{
		employeeRouter.GET("", employeeHandler.List)
		employeeRouter.POST("", middleware.APIKey(config.Passphrase), employeeHandler.Add)
		employeeRouter.GET("/:id", employeeHandler.Detail)
	}

	payrollRouter := router.Group("/payrolls")
	{
		payrollRouter.GET("", payrollHandler.List)
		payrollRouter.POST("", middleware.APIKey(config.Passphrase), payrollHandler.Add)
		payrollRouter.GET("/:id", payrollHandler.Detail)
	}

	salaryRouter := router.Group("/salaries")
	{
		salaryRouter.GET("", salaryHandler.List)
	}

	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if srvErr := server.Run(srv, logger); srvErr != nil {
		logger.Fatal().Err(srvErr).Msg("server shutdown failed")
	}
}
