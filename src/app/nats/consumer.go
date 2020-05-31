package nats

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
	"time"
)

type Consumer interface {
	Next(timeout time.Duration) []byte
	Subscribe(subject string) error
	SubscribeAsync(subject string) (chan *nats.Msg, error)
	Reply(subject string)
	Close()
}

type consumerImpl struct {
	c *nats.Conn
	s *nats.Subscription
}

var consumer Consumer

func NewConsumer() Consumer {
	var once sync.Once
	once.Do(func() {
		conn := Connect("consumer", nats.DefaultURL)
		consumer = &consumerImpl{
			c: conn,
		}
	})
	return consumer
}

func (r *consumerImpl) Subscribe(subject string) error {
	sub, err := r.c.SubscribeSync(subject)
	if err != nil {
		log.Println(fmt.Sprintf("Error when subscribe: %s", subject))
		return err
	}
	r.s = sub
	return nil
}

func (r *consumerImpl) Next(timeout time.Duration) []byte {
	var data []byte
	m, err := r.s.NextMsg(timeout)
	if err != nil {
		log.Println("Error when get next message\n" + err.Error())
		return data
	}

	return m.Data
}

func (r *consumerImpl) SubscribeAsync(subject string) (chan *nats.Msg, error) {
	ch := make(chan *nats.Msg, 64)
	_, err := r.c.ChanSubscribe(subject, ch)
	if err != nil {
		return nil, errors.New("Error when try create subscribe")
	}
	return ch, nil
}

func (r *consumerImpl) Reply(sub string) {
	var pause bool
	for !pause {
		r.c.Subscribe(sub, func(m *nats.Msg) {
			fmt.Println("Request: " + string(m.Data))
			r.c.Publish(m.Reply, []byte("pong"))
			pause = true
		})
		time.Sleep(1 * time.Second)
	}
}

func (r *consumerImpl) Close() {
	r.c.Close()
}
