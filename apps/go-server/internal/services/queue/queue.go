// Package queue wraps RabbitMQ so slow I/O (email, reports, bulk jobs) never
// runs on the API request path — the API publishes, the worker (cmd/worker)
// consumes. Keeping RabbitMQ (not Azure Service Bus) per the rewrite's
// messaging decision: it's already running locally and the code path matches
// the existing FastAPI monolith's services/queue.py.
package queue

import (
	"context"
	"encoding/json"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher lazily dials RabbitMQ and re-dials if the connection drops.
type Publisher struct {
	url string

	mu   sync.Mutex
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewPublisher builds a Publisher for the given AMQP URL. No connection is
// made until the first Publish call.
func NewPublisher(url string) *Publisher {
	return &Publisher{url: url}
}

func (p *Publisher) channel() (*amqp.Channel, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn != nil && !p.conn.IsClosed() && p.ch != nil {
		return p.ch, nil
	}
	conn, err := amqp.Dial(p.url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	p.conn, p.ch = conn, ch
	return ch, nil
}

// Publish declares queueName durable and publishes body as a persistent JSON
// message, matching services/queue.py's publish_message behaviour.
func (p *Publisher) Publish(ctx context.Context, queueName string, body any) error {
	ch, err := p.channel()
	if err != nil {
		return err
	}
	if _, err := ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		return err
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return ch.PublishWithContext(ctx, "", queueName, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         payload,
	})
}

// Close releases the channel and connection, if open.
func (p *Publisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.ch != nil {
		_ = p.ch.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

// Consume dials a fresh connection and returns a delivery channel for
// queueName. Intended for long-lived worker processes (one call at startup).
func Consume(url, queueName string) (<-chan amqp.Delivery, *amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, nil, err
	}
	if err := ch.Qos(10, 0, false); err != nil {
		_ = conn.Close()
		return nil, nil, err
	}
	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		_ = conn.Close()
		return nil, nil, err
	}
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		_ = conn.Close()
		return nil, nil, err
	}
	return msgs, conn, nil
}
