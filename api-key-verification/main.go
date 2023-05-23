package main

import (
	"api-key-verification/config"
	"api-key-verification/handler"
	"api-key-verification/helper/crypt"
	"api-key-verification/helper/logger"
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

	crypt := crypt.New(config.Passphrase)

	pingHandler := handler.NewPingGinHandler()
	apiKeyHandler := handler.NewAPIKeyGinHandler(crypt)

	router := gin.New()

	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())
	// router.Use(middleware.APIKey(config.Passphrase))

	router.GET("/ping", pingHandler.Ping)

	router.GET("/verify", apiKeyHandler.Verify)
	router.POST("/generate", apiKeyHandler.Generate)

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
