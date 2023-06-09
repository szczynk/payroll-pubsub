//nolint:dupl // this is why
package repository

import (
	"encoding/json"
	"fundamental-payroll-gin/config/rmq"
	"fundamental-payroll-gin/model"

	"github.com/rs/zerolog"
)

type EmployeeRMQPubRepo struct {
	rmq           rmq.InterfaceRMQ
	logger        *zerolog.Logger
	exchangeName  string
	routingKeyPub string
	routingKeyCon string
}

func NewEmployeeRMQPubRepo(rmq rmq.InterfaceRMQ, logger *zerolog.Logger) EmployeeRMQPubRepoI {
	r := new(EmployeeRMQPubRepo)
	r.rmq = rmq
	r.logger = logger
	r.exchangeName = "payrollapp.employee"
	r.routingKeyPub = "employee.insert"
	r.routingKeyCon = "employee.created"
	return r
}

func (r *EmployeeRMQPubRepo) Add(employee *model.Employee) (*model.Employee, error) {
	employeeBytes, errJSON := json.Marshal(employee)
	if errJSON != nil {
		return nil, errJSON
	}

	errPub := r.rmq.Publish(r.exchangeName, r.routingKeyPub, employeeBytes)
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

	done := make(chan *model.Employee, 1)
	go func() {
		for msg := range msgs {
			r.logger.Debug().Msgf("from consumer-payroll Received a message: %s", msg.Body)

			newEmployee := new(model.Employee)

			errJSONUn := json.Unmarshal(msg.Body, &newEmployee)
			if errJSONUn != nil {
				r.logger.Error().Err(errJSONUn).Msg("json.Unmarshal err")
				_ = msg.Reject(false)
			}

			// r.logger.Debug().Msgf("encoded new employee: %v", newEmployee)
			done <- newEmployee
			_ = msg.Ack(false)
		}
	}()

	return <-done, nil
}
