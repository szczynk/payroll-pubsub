package main

import (
	"consumer-payroll/config"
	"consumer-payroll/config/db"
	"consumer-payroll/config/rmq"
	"consumer-payroll/handler"
	"consumer-payroll/helper/logger"
	"consumer-payroll/repository"
	"consumer-payroll/usecase"
	"log"
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
	payrollRepo := repository.NewPayrollGormRepository(gormDB.DB)

	employeeUC := usecase.NewEmployeeUsecase(employeeRepo)
	payrollUC := usecase.NewPayrollUsecase(payrollRepo)

	employeeHandler := handler.NewEmployeeRMQHandler(rmq, logger, employeeUC)
	payrollHandler := handler.NewPayrollRMQHandler(rmq, logger, payrollUC)

	var forever chan struct{}

	go employeeHandler.Add()
	go payrollHandler.Add()

	logger.Debug().Msg("[*] To exit press CTRL+C")
	<-forever
}
