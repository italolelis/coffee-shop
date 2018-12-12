package coffee

import (
	"context"
	"time"
)

type Coffee interface {
	Brew(context.Context)
	Match(string) bool
}

type Espresso struct{}

func (e *Espresso) Brew(ctx context.Context) {
	time.Sleep(3 * time.Second)
}

func (e *Espresso) Match(coffeeType string) bool {
	return coffeeType == "espresso"
}

type Cappuccino struct{}

func (c *Cappuccino) Brew(ctx context.Context) {
	time.Sleep(6 * time.Second)
}

func (c *Cappuccino) Match(coffeeType string) bool {
	return coffeeType == "cappuccino"
}

type Latte struct{}

func (l *Latte) Brew(ctx context.Context) {
	time.Sleep(8 * time.Second)
}

func (l *Latte) Match(coffeeType string) bool {
	return coffeeType == "latte"
}
