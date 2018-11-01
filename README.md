# Coffee Shop

This is a example of how we can design an application that is ready to handle unexpected events.

## Reception

Reception is the service that is receiving new orders and sending to the `barista` to be done.

### Domain

* We have a simple `orders` API where you can place your orders. An order is only placed after the payment is done successfuly.
* To keep things simple, we mock the `payments` service with [wiremock](http://wiremock.org) in a separate container where we can simulate failure when necessary.

## Barista

Barista is the service that actually prepares your coffee and make sure you can have the best experience possible.

The `barista` will get an order request from rabbitmq and start preparing your coffee. Once the coffee is done another message is published
and whoever want's to interact with that will be able to.
There is a read API to fetch all `orders` that are ready to be taken
