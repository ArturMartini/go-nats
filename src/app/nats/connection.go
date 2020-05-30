package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func Connect(context, url string) *nats.Conn {
	c, err := nats.Connect(url, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(fmt.Sprintf("Connected %s in NATS %s", context, url))
	return c
}
