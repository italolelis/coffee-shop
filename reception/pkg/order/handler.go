package order

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/golang/protobuf/proto"
	"github.com/italolelis/kit/log"
	"github.com/italolelis/kit/proto/order"
	"github.com/italolelis/reception/pkg/coffees"
	"github.com/rafaeljesus/rabbus"
	"github.com/satori/go.uuid"
	"go.opencensus.io/trace"
)

const (
	exchangeName = "orders"
	exchangeType = "topic"
)

type createOrderRequest struct {
	Name  string `json:"name"`
	Items []struct {
		Type string `json:"type"`
		Size string `json:"size"`
	} `json:"items"`
}

type Handler struct {
	wRepo          WriteRepository
	rRepo          ReadRepository
	coffeeReadRepo coffees.ReadRepository
	es             *rabbus.Rabbus
}

// NewHandler creates a new instance of OrderHandler
func NewHandler(wRepo WriteRepository, rRepo ReadRepository, coffeeReadRepo coffees.ReadRepository, es *rabbus.Rabbus) *Handler {
	return &Handler{wRepo: wRepo, rRepo: rRepo, coffeeReadRepo: coffeeReadRepo, es: es}
}

// CreateOrder is the handler for order creation
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq createOrderRequest

	ctx := r.Context()
	logger := log.WithContext(ctx)
	ctx, span := trace.StartSpan(ctx, "handler:create.order", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&orderReq); err != nil {
		http.Error(w, "could not parse the request body", http.StatusBadRequest)
		return
	}

	o := NewOrder(orderReq.Name)
	for _, i := range orderReq.Items {
		coffee, err := h.coffeeReadRepo.FindOneByName(ctx, i.Type)
		if err != nil {
			logger.With("err", err, "name", i.Type).Error("could not find the coffee")
			http.Error(w, fmt.Sprintf("unfortunately we don't have %s", i.Type), http.StatusBadRequest)
			return
		}

		o.Items = append(o.Items, &Item{
			Coffee: coffee,
			Size:   i.Size,
		})
	}

	logger.Debugw("creating order", "order_id", o.ID)
	if err := h.wRepo.Add(ctx, o); err != nil {
		http.Error(w, "could not save your order", http.StatusInternalServerError)
		return
	}

	ev := order.Created{
		ID:           o.ID.String(),
		CustomerName: o.CustomerName,
	}
	for _, i := range o.Items {
		ev.Items = append(ev.Items, &order.OrderItem{Type: i.Coffee.Name, Size: i.Size})
	}

	if err := h.sendEvent(ctx, "orders.created", &ev); err != nil {
		http.Error(w, "could not send an event of your order", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("/orders/%s", o.ID.String()))
	w.WriteHeader(http.StatusCreated)
}

// GetOrder is the handler for getting a single order
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.WithContext(ctx)
	ctx, span := trace.StartSpan(ctx, "handler:get.order", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		logger.Errorw("invalid id", "order_id", chi.URLParam(r, "id"))
		http.Error(w, "invalid order id provided", http.StatusBadRequest)
		return
	}

	data, err := h.rRepo.FindOneByID(ctx, id)
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

func (h *Handler) sendEvent(ctx context.Context, key string, payload proto.Message) error {
	ctx, span := trace.StartSpan(ctx, "rabbitmq:send", trace.WithSpanKind(trace.SpanKindClient))
	span.AddAttributes(trace.StringAttribute("exchange", exchangeName), trace.StringAttribute("key", key))
	defer span.End()

	data, err := proto.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal the message before sending to the message stream: %s", err.Error())
	}

	msg := rabbus.Message{
		Exchange: exchangeName,
		Kind:     exchangeType,
		Key:      key,
		Payload:  data,
	}

	h.es.EmitAsync() <- msg

	return nil
}
