package main

import (
	"fundamental-payroll-gin/config"
	"fundamental-payroll-gin/config/db"
	"fundamental-payroll-gin/config/rmq"
	"fundamental-payroll-gin/handler"
	"fundamental-payroll-gin/helper/logger"
	"fundamental-payroll-gin/repository"
	"fundamental-payroll-gin/server"
	"fundamental-payroll-gin/usecase"
	"log"
	"net/http"
	"time"
)

func main() {
	config, confErr := config.LoadConfig()
	if confErr != nil {
		log.Fatalf("load config err:%s", confErr)
	}

	logger := logger.New(config.Debug)

	rmq, errRmq := rmq.NewRMQ(config.RabbitMQ.URL)
	if errRmq != nil {
		logger.Fatal().Err(errRmq).Msg("rabbitmq failed to connect")
	}
	logger.Debug().Msg("rabbitmq connected")
	rmq.SetupBlockingNotifications(logger)

	gormDB, errDB := db.NewGormDB(config.Debug, config.Database.Driver, config.Database.URL)
	if errDB != nil {
		logger.Fatal().Err(errDB).Msg("gorm postgres failed to connect")
	}
	logger.Debug().Msg("gorm postgres connected")

	defer func() {
		logger.Debug().Msg("closing rabbitmq")
		_ = rmq.ConnClose()
		_ = rmq.ChClose()

		logger.Debug().Msg("closing gorm postgres")
		_ = gormDB.Close()
	}()

	employeeRepo := repository.NewEmployeeGormRepository(gormDB.DB)
	employeePub := repository.NewEmployeeRMQPubRepo(rmq, logger)

	payrollRepo := repository.NewPayrollGormRepository(gormDB.DB)
	payrollPub := repository.NewPayrollRMQPubRepo(rmq, logger)
	salaryRepo := repository.NewSalaryGormRepository(gormDB.DB)

	employeeUC := usecase.NewEmployeeUsecase(employeeRepo, employeePub, salaryRepo)
	payrollUC := usecase.NewPayrollUsecase(payrollRepo, payrollPub, employeeRepo, salaryRepo)
	salaryUC := usecase.NewSalaryUsecase(salaryRepo)

	employeeHandler := handler.NewEmployeeGinHandler(employeeUC)
	payrollHandler := handler.NewPayrollGinHandler(payrollUC)
	salaryHandler := handler.NewSalaryGinHandler(salaryUC)

	router := server.NewRouter(config.Debug, logger, employeeHandler, payrollHandler, salaryHandler)

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
