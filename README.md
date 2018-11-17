# Coffee Shop

This is a example of how we can design an application that is ready to handle unexpected events.

## Getting Started

You can easily spin up the whole structure by executing:

```sh
docker-compose up -d
```

This will setup our 3 services and will make sure they are talking to each other. If you want to play around with it and simluate failure you can run the [failure scenarios](#simulating-failure).

## Reception Service

Reception is the service that is receiving new orders and sending to the `barista` to be done.

* We have a simple `orders` API where you can place your orders. An order is only placed after the payment is done successfuly.
* To keep things simple, we mock the `payments` service with [wiremock](http://wiremock.org) in a separate container where we can simulate failure when necessary.

## Barista Service

Barista is the service that actually prepares your coffee and make sure you can have the best experience possible.

The `barista` will get an order request from rabbitmq and start preparing your coffee. Once the coffee is done another message is published
and whoever want's to interact with that will be able to.
There is a read API to fetch all `orders` that are ready to be taken

## Simulating Failure

### Breaking the message broker

```sh
docker-compose stop rabbitmq
```

```sh
make orders
```
