package handler

import (
	"consumer-payroll/config/rmq"
	"consumer-payroll/model"
	"consumer-payroll/usecase"
	"encoding/json"

	"github.com/rs/zerolog"
)

type PayrollRMQHandler struct {
	rmq           rmq.InterfaceRMQ
	payrollUC     usecase.PayrollUsecaseI
	logger        *zerolog.Logger
	exchangeName  string
	routingKeyPub string
	routingKeyCon string
}

func NewPayrollRMQHandler(
	rmq rmq.InterfaceRMQ,
	logger *zerolog.Logger,
	payrollUC usecase.PayrollUsecaseI,
) PayrollRMQHandlerI {
	h := new(PayrollRMQHandler)
	h.rmq = rmq
	h.payrollUC = payrollUC
	h.logger = logger
	h.exchangeName = "payrollapp.payroll"
	h.routingKeyPub = "payroll.created"
	h.routingKeyCon = "payroll.insert"
	return h
}

func (h *PayrollRMQHandler) Add() {
	msgs, errSub := h.rmq.Subscribe(h.exchangeName, h.routingKeyCon, false)
	if errSub != nil {
		h.logger.Error().Err(errSub).Msg("rmq.Subscribe err")
		return
	}

	payroll := new(model.Payroll)
	for msg := range msgs {
		h.logger.Debug().Msgf("from fundamental-payroll-gin Received a message: %s", msg.Body)

		errJSONUn := json.Unmarshal(msg.Body, &payroll)
		if errJSONUn != nil {
			h.logger.Error().Err(errJSONUn).Msg("json.Unmarshal err")
			// Assume the message is rejected due to an invalid format
			// The delivery variable represents the received message delivery
			// Reject the message
			_ = msg.Reject(false)
			continue
		}

		newPayroll, ucErr := h.payrollUC.Add(payroll)
		if ucErr != nil {
			h.logger.Error().Err(ucErr).Msg("payrollUC.Add err")
			// Assume an error occurred while processing the message
			// The delivery variable represents the received message delivery
			// Negatively acknowledge the message
			_ = msg.Nack(false, false)
			continue
		}

		payrollBytes, errJSON := json.Marshal(newPayroll)
		if errJSON != nil {
			h.logger.Error().Err(errJSON).Msg("json.Marshal err")
			_ = msg.Nack(false, false)
			continue
		}

		errPub := h.rmq.Publish(h.exchangeName, h.routingKeyPub, payrollBytes)
		if errPub != nil {
			h.logger.Error().Err(errPub).Msg("rmq.Publish err")
			_ = msg.Nack(false, false)
			continue
		}

		h.logger.Debug().Msgf("rmq.PubSub on Add success with id:%v", newPayroll.ID)

		// Assume you have received a message and processed it successfully
		// The delivery variable represents the received message delivery
		// Acknowledge the message
		_ = msg.Ack(false)
	}
}
