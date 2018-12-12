package staff

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/italolelis/barista/pkg/coffee"
	"github.com/italolelis/kit/log"
	"github.com/italolelis/kit/proto/order"
	"github.com/rafaeljesus/rabbus"
)

type OrderDone struct {
	CustomerName string
	DoneBy       *Barista
	Type         string
	Size         string
}

type Barista struct {
	Name   string
	Skills []coffee.Coffee
}

func (b *Barista) Prepare(ctx context.Context, messages chan rabbus.ConsumerMessage, result chan<- *OrderDone) {
	logger := log.WithContext(ctx)

	for {
		m, ok := <-messages
		if !ok {
			logger.Debug("stop listening messages!")
			break
		}

		o := order.Created{}
		err := proto.Unmarshal(m.Body, &o)
		if err != nil {
			logger.Errorw("unmarshal error", "err", err)
		}

		logger.Infow("preparing order", "id", o.ID)
		for _, i := range o.Items {
			for _, s := range b.Skills {
				if s.Match(i.Type) {
					s.Brew(ctx)
					result <- &OrderDone{
						CustomerName: o.CustomerName,
						DoneBy:       b,
						Type:         i.Type,
						Size:         i.Size,
					}
				}
			}
		}

		if err := m.Ack(false); err != nil {
			logger.Errorw("failed to ack the message", "err", err)
		}

		logger.Debug("message was consumed")
	}
}
