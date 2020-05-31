package nats

import (
	nats "github.com/nats-io/nats.go"
	"log"
	"sync"
	"time"
)

type Producer interface {
	Send(string, []byte) error
	Request(string, []byte, time.Duration) string
	Close()
}

type producerImpl struct {
	conn *nats.Conn
}

var producer Producer

func NewProducer() Producer {
	var once sync.Once
	once.Do(func() {
		producer = &producerImpl{
			conn: Connect("producer", nats.DefaultURL),
		}
	})
	return producer
}

func (r producerImpl) Send(subject string, bytes []byte) error {
	err := r.conn.Publish(subject, bytes)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (r producerImpl) Request(subject string, bytes []byte, timeout time.Duration) string {
	start := time.Now()
	for time.Now().Before(start.Add(timeout)) {
		msg, err := r.conn.Request(subject, bytes, 1 * time.Second)
		if err != nil {
			continue
		}
		return string(msg.Data)
	}
	return ""
}

func (r producerImpl) Close() {
	r.conn.Close()
}
