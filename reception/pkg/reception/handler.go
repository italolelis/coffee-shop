package reception

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/italolelis/kit/log"
	"github.com/italolelis/kit/proto/order"
	"github.com/satori/go.uuid"
	"github.com/rafaeljesus/rabbus"
	"github.com/golang/protobuf/proto"
)

const (
	exchangeName = "orders"
	exchangeType = "topic"
)

// CreateOrderHandler is the hander for order creation
func CreateOrderHandler(repo WriteRepository, eventStream *rabbus.Rabbus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := log.WithContext(ctx)

		o := NewOrder()
		d := json.NewDecoder(r.Body)
		if err := d.Decode(&o); err != nil {
			http.Error(w, "could not parse the request body", http.StatusBadRequest)
			return
		}

		logger.Debugw("creating order", "order_id", o.ID)
		if err := repo.Add(ctx, o); err != nil {
			http.Error(w, "could not save your order", http.StatusInternalServerError)
			return
		}

		ev := order.Created{
			ID: o.ID.String(),
		}
		for _, i := range o.Items {
			ev.Items = append(ev.Items, i.Type)
		}
		
		if err:= sendEvent(eventStream, "orders.created", &ev); err != nil {
			http.Error(w, "could not send an event of your order", http.StatusInternalServerError)
			return
		}

		w.Header().Add("Location", fmt.Sprintf("/orders/%s", o.ID.String()))
		w.WriteHeader(http.StatusCreated)
	}
}

// GetOrderHandler is the hander for getting a single order
func GetOrderHandler(repo ReadRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := log.WithContext(ctx)

		id, err := uuid.FromString(chi.URLParam(r, "id"))
		if err != nil {
			logger.Errorw("invalid id", "order_id", chi.URLParam(r, "id"))
			http.Error(w, "invalid id provided", http.StatusBadRequest)
			return
		}

		data, err := repo.FindOneByID(ctx, id)
		if err != nil {
			logger.Errorw("could not get the order", "order_id", id, "err", err)
			if err == ErrOrderNotFound {
				http.Error(w, "invalid order id", http.StatusNotFound)
				return
			}

			http.Error(w, "unkown error when retreiving the order", http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, data)
	}
}

func sendEvent(r *rabbus.Rabbus, key string, payload proto.Message) error {
	data, err := proto.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal the message before sending to the message stream: %s", err.Error())
	}

	msg := rabbus.Message{
		Exchange: exchangeName,
		Kind:    exchangeType,
		Key:     key,
		Payload:  data,
	}

	r.EmitAsync() <- msg

	return nil
}
