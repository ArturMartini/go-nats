package nats

import (
	nats "github.com/nats-io/nats.go"
	"log"
	"sync"
)


type Producer interface {
	Send(string, []byte) error
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