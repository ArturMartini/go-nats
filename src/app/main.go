package main

import (
	"app/nats"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	//This logic is necessary because NATS for default is are in memory state
	//Im just create two parallel routines and force consumer and publish at same time
	subject := "hello"
	str := "ping"
	wg.Add(2)
	go consumerSync(subject)
	go producer(subject, str)
	wg.Wait()
}

func producer(subject, message string) {
	time.Sleep(2 * time.Second)
	prod := nats.NewProducer()
	prod.Send(subject, []byte(message))
	wg.Done()
}

func consumerSync(subject string) {
	cons := nats.NewConsumer(subject)
	message := cons.Next(10 * time.Second)
	fmt.Println(string(message))
	wg.Done()
}
