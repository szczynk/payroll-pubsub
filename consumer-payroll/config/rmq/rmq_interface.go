//go:generate mockery --output=../../mocks --name InterfaceRMQ
package rmq

import (
	"github.com/rabbitmq/amqp091-go"
)

type InterfaceRMQ interface {
	ConnClose() error
	ChClose() error
	Publish(exchangeName, routingKeyPub string, dataBytes []byte) error
	Subscribe(exchangeName, routingKeyCon string, autoAck bool) (<-chan amqp091.Delivery, error)
}
