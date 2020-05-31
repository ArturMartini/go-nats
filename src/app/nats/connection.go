package nats

import (
	"github.com/nats-io/nats.go"
	"log"
)

func Connect(context, url string) *nats.Conn {
	c, err := nats.Connect(url, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	return c
}
