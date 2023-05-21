//nolint:dupl // this is why
package repository

import (
	"encoding/json"
	"fundamental-payroll-gin/config/rmq"
	"fundamental-payroll-gin/model"

	"github.com/rs/zerolog"
)

type PayrollRMQPubRepo struct {
	rmq           *rmq.RMQ
	logger        *zerolog.Logger
	exchangeName  string
	routingKeyPub string
	routingKeyCon string
}

func NewPayrollRMQPubRepo(rmq *rmq.RMQ, logger *zerolog.Logger) PayrollRMQPubRepoI {
	r := new(PayrollRMQPubRepo)
	r.rmq = rmq
	r.logger = logger
	r.exchangeName = "payrollapp.payroll"
	r.routingKeyPub = "payroll.insert"
	r.routingKeyCon = "payroll.created"
	return r
}

func (r *PayrollRMQPubRepo) Add(payroll *model.Payroll) (*model.Payroll, error) {
	payrollBytes, errJSON := json.Marshal(payroll)
	if errJSON != nil {
		return nil, errJSON
	}

	errPub := r.rmq.Publish(r.exchangeName, r.routingKeyPub, payrollBytes)
	if errPub != nil {
		return nil, errPub
	}

	// note: tldr
	// the consumer keep stuck whether
	// context deadline exceed or closed connection/channel or just keep circle in around
	// under conditions:
	// 1. autoAck false regardless d.Ack exist or not
	// 2. d.Ack not exist regardless autoAck true or not
	// trust me bro
	// ! keep this intact
	msgs, err := r.rmq.Subscribe(r.exchangeName, r.routingKeyCon, true)
	if err != nil {
		return nil, err
	}

	done := make(chan *model.Payroll)
	go func() {
		for msg := range msgs {
			r.logger.Debug().Msgf("from consumer-payroll Received a message: %s", msg.Body)

			newPayroll := new(model.Payroll)

			errJSONUn := json.Unmarshal(msg.Body, &newPayroll)
			if errJSONUn != nil {
				r.logger.Error().Err(errJSONUn).Msg("json.Unmarshal err")
				_ = msg.Reject(false)
			}

			// r.logger.Debug().Msgf("encoded new payroll: %v", newPayroll)
			done <- newPayroll
			_ = msg.Ack(false)
		}
	}()

	return <-done, nil
}
