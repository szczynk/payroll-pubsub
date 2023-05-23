package handler

import (
	"consumer-payroll/config/rmq"
	"consumer-payroll/model"
	"consumer-payroll/usecase"
	"encoding/json"

	"github.com/rs/zerolog"
)

type EmployeeRMQHandler struct {
	rmq           rmq.InterfaceRMQ
	employeeUC    usecase.EmployeeUsecaseI
	logger        *zerolog.Logger
	exchangeName  string
	routingKeyPub string
	routingKeyCon string
}

func NewEmployeeRMQHandler(
	rmq rmq.InterfaceRMQ,
	logger *zerolog.Logger,
	employeeUC usecase.EmployeeUsecaseI,
) EmployeeRMQHandlerI {
	h := new(EmployeeRMQHandler)
	h.rmq = rmq
	h.employeeUC = employeeUC
	h.logger = logger
	h.exchangeName = "payrollapp.employee"
	h.routingKeyPub = "employee.created"
	h.routingKeyCon = "employee.insert"
	return h
}

func (h *EmployeeRMQHandler) Add() {
	msgs, errSub := h.rmq.Subscribe(h.exchangeName, h.routingKeyCon, false)
	if errSub != nil {
		h.logger.Error().Err(errSub).Msg("rmq.Subscribe err")
		return
	}

	employee := new(model.Employee)
	for msg := range msgs {
		h.logger.Debug().Msgf("from fundamental-payroll-gin Received a message: %s", msg.Body)

		errJSONUn := json.Unmarshal(msg.Body, &employee)
		if errJSONUn != nil {
			h.logger.Error().Err(errJSONUn).Msg("json.Unmarshal err")
			// Assume the message is rejected due to an invalid format
			// The delivery variable represents the received message delivery
			// Reject the message
			_ = msg.Reject(false)
			continue
		}

		newEmployee, ucErr := h.employeeUC.Add(employee)
		if ucErr != nil {
			h.logger.Error().Err(ucErr).Msg("employeeUC.Add err")
			// Assume an error occurred while processing the message
			// The delivery variable represents the received message delivery
			// Negatively acknowledge the message
			_ = msg.Nack(false, false)
			continue
		}

		employeeBytes, errJSON := json.Marshal(newEmployee)
		if errJSON != nil {
			h.logger.Error().Err(errJSON).Msg("json.Marshal err")
			_ = msg.Nack(false, false)
			continue
		}

		errPub := h.rmq.Publish(h.exchangeName, h.routingKeyPub, employeeBytes)
		if errPub != nil {
			h.logger.Error().Err(errPub).Msg("rmq.Publish err")
			_ = msg.Nack(false, false)
			continue
		}

		// Assume you have received a message and processed it successfully
		// The delivery variable represents the received message delivery
		// Acknowledge the message
		_ = msg.Ack(false)
		h.logger.Debug().Msgf("rmq.PubSub on Add success with id:%v", newEmployee.ID)
	}
}
