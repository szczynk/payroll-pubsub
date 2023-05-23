package rmq

import (
	"errors"
	"fmt"
	"fundamental-payroll-gin/helper/timeout"

	"github.com/rabbitmq/amqp091-go"
)

type RMQ struct {
	url  string
	Conn *amqp091.Connection
	Ch   *amqp091.Channel
}

func NewRMQ(url string) (InterfaceRMQ, error) {
	if url == "" {
		return nil, errors.New("no database url")
	}

	rmq := new(RMQ)
	err := rmq.init(url)
	if err != nil {
		return nil, err
	}

	return rmq, nil
}

func (r *RMQ) init(url string) error {
	r.url = url

	conn, errConn := amqp091.Dial(r.url)
	if errConn != nil {
		return fmt.Errorf("failed to connect to rabbitmq: %w", errConn)
	}

	ch, errCh := conn.Channel()
	if errCh != nil {
		return fmt.Errorf("failed to open a channel: %w", errCh)
	}

	r.Conn = conn
	r.Ch = ch
	return nil
}

func (r *RMQ) ensureConnectionAndChannel() error {
	if r.Conn == nil || r.Conn.IsClosed() {
		conn, errConn := amqp091.Dial(r.url)
		if errConn != nil {
			return fmt.Errorf("failed to reconnect to rabbitmq: %w", errConn)
		}
		r.Conn = conn
	}

	if r.Ch == nil || r.Ch.IsClosed() {
		ch, errCh := r.Conn.Channel()
		if errCh != nil {
			return fmt.Errorf("failed to reopen a channel: %w", errCh)
		}
		r.Ch = ch
	}

	return nil
}

func (r *RMQ) ConnClose() error {
	return r.Conn.Close()
}

func (r *RMQ) ChClose() error {
	return r.Ch.Close()
}

func (r *RMQ) Publish(exchangeName, routingKeyPub string, dataBytes []byte) error {
	errEnsure := r.ensureConnectionAndChannel()
	if errEnsure != nil {
		return errEnsure
	}

	errExc := r.Ch.ExchangeDeclare(
		exchangeName, // exchange name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if errExc != nil {
		return errExc
	}

	// ctx:=context.Background()
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	errPub := r.Ch.PublishWithContext(
		ctx,
		exchangeName,  // exchange name
		routingKeyPub, // routing key
		false,         // mandatory
		false,         // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        dataBytes,
		})
	if errPub != nil {
		return errPub
	}

	return nil
}

func (r *RMQ) Subscribe(exchangeName, routingKeyCon string, autoAck bool) (<-chan amqp091.Delivery, error) {
	errEnsure := r.ensureConnectionAndChannel()
	if errEnsure != nil {
		return nil, errEnsure
	}

	errExc := r.Ch.ExchangeDeclare(
		exchangeName, // exchange name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if errExc != nil {
		return nil, errExc
	}

	q, errQ := r.Ch.QueueDeclare(
		routingKeyCon+"q", // queue name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if errQ != nil {
		return nil, errQ
	}

	errQB := r.Ch.QueueBind(
		q.Name,        // queue name
		routingKeyCon, // routing key
		exchangeName,  // exchange
		false,         // no-wait
		nil,           // arguments
	)
	if errQB != nil {
		return nil, errQB
	}

	errQos := r.Ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if errQos != nil {
		return nil, errQos
	}

	msgs, errCon := r.Ch.Consume(
		q.Name,  // queue name
		"",      // consumer tag
		autoAck, // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if errCon != nil {
		return nil, errCon
	}
	return msgs, nil
}
