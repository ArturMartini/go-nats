package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
	"time"
)

type Consumer interface {
	Next(timeout time.Duration) []byte
}

type consumerImpl struct {
	c *nats.Conn
	s *nats.Subscription
}

var consumer Consumer

func NewConsumer(subject string) Consumer {
	var once sync.Once
	once.Do(func() {
		conn := Connect("consumer", nats.DefaultURL)
		sub, err := conn.SubscribeSync(subject)
		if err != nil {
			log.Fatal(fmt.Sprintf("Error when subscribe: %s", subject))
		}

		consumer = &consumerImpl{
			c: conn,
			s: sub,
		}
	})
	return consumer
}

func (r consumerImpl) Next(timeout time.Duration) []byte {
	var data []byte
	m, err := r.s.NextMsg(timeout)
	if err != nil {
		log.Println("Error when get next message\n" + err.Error())
		return data
	}

	return m.Data
}
