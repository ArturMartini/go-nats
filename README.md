# go-nats 
This project have an examples of how to use client integration with NATS

## For run project you need
* Go installed
* Docker running in local machine
* Run docker pull nats
* Run docker run -d --name nats-main -p 4222:4222 -p 6222:6222 -p 8222:8222 nats
* Run go build main.go
* Run ./main.go

### We implemented the scenarios below
* Producer messages
* Consumer Sync messages
* Consumer Async messages
* Producer with Request and wait reply message

#### If run succesfully you see messages below
* Sync: ping
* Async: 1
* Async: 2
* Async: 3
* Async: 4
* Async: 5
* Request: ping
* Reply: pong

##### Note. In this example, using NATS in memory, So some wait groups and sleeps are implemented for produce and consume simulation at same time


