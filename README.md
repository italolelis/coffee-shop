# Coffee Shop

[![Build Status](https://travis-ci.com/italolelis/coffee-shop.svg?branch=master)](https://travis-ci.com/italolelis/coffee-shop)

This is a example of how we can design an application that is ready to handle unexpected events.

## Getting Started

You can easily spin up the whole system by executing:

```sh
docker-compose up -d
```

This will setup our 2 services and will make sure they are talking to each other. 
If you want to play around with it and simulate failure you can run the [failure scenarios](#simulating-failure).

## Reception Service

Think about a coffee-shop. The first thing you normally do when you enter one is to go to the reception 
and order a coffee. Reception is the service that is getting new orders and sending to the `barista` to be done.

In this service we have:
* RabbitMQ to send messages that are send with [protocol buffers](/configs/proto)
* Tracing and monitoring using [open census](https://github.com/census-instrumentation/opencensus-go)
* All calls to external dependencies are wrapped around a Circuit Breaker. You can use a Hystrix dashboard to check for the circuits.
* To keep things simple, we mock the `payments` service with [wiremock](http://wiremock.org) in a separate container where we can simulate failure when necessary.

## Barista Service

Barista is the service that actually prepares your coffee and make sure you can have the best experience possible.

The `barista` will get an order request from RabbitMQ and start preparing your coffee. Once the coffee is ready 
another message is published and whoever want's to interact with that will be able to.

## Tools

* Jeager UI: http://localhost:16686
* Prometheus: 
* Hystrix Dashboard: 

## Simulating Failure

### Breaking the message broker

```sh
docker-compose stop rabbitmq
```

```sh
make orders
```
