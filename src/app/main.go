package main

import (
	"app/nats"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup
var subject string = "hello"
var str = "ping"

func main() {
	scenarioProducerAndConsumerSync()
	scenarioProducerAndConsumerAsync()
	scenarioRequestReply()
}

func producer(subject, message string) {
	time.Sleep(2 * time.Second)
	prod := nats.NewProducer()
	defer prod.Close()
	prod.Send(subject, []byte(message))
}

func consumerSync(subject string) {
	consumer := nats.NewConsumer()
	consumer.Subscribe(subject)
	defer consumer.Close()
	message := consumer.Next(5 * time.Second)
	fmt.Println(string(message))
	wg.Done()
}

func consumerAsync(subject string) {
	consumer := nats.NewConsumer()
	defer consumer.Close()
	ch, err := consumer.SubscribeAsync(subject)
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		select {
		case msg := <-ch:
			fmt.Println(string(msg.Data))
		}
	}
}

func producerPerSecond(subject, message string, seconds int) {
	prod := nats.NewProducer()
	defer prod.Close()
	for i := 0; i <= seconds; i++ {
		prod.Send(subject, []byte(strconv.Itoa(i)))
		time.Sleep(1 * time.Second)
	}
	wg.Done()
}

func scenarioProducerAndConsumerSync() {
	//This logic is necessary because NATS for default is are in memory state
	//Im just create two parallel routines and force consumer and publish at same time
	fmt.Println("Execute Producer and Consumer Sync")
	wg.Add(1)
	go consumerSync(subject)
	go producer(subject, str)
	wg.Wait()
	fmt.Println("Execute Producer and Consumer Sync finish")
}

func scenarioProducerAndConsumerAsync() {
	fmt.Println("Execute Producer and Consumer Async")
	wg.Add(1)
	go consumerAsync(subject)
	go producerPerSecond(subject, "ping-async", 5)
	wg.Wait()
	fmt.Println("Execute Producer and Consumer Async finish")
}

func scenarioRequestReply() {
	subject := "request-reply"
	wg.Add(2)
	go producerWithRequest(subject)
	go consumerWithReply(subject)
	wg.Wait()
}

func consumerWithReply(subject string) {
	consumer := nats.NewConsumer()
	defer consumer.Close()
	consumer.Reply(subject)
	wg.Done()
}

func producerWithRequest(subject string) {
	prod := nats.NewProducer()
	defer prod.Close()
	reply := prod.Request(subject, []byte("ping"), 5 *time.Second)
	fmt.Println("Reply: " + reply)
	wg.Done()
}