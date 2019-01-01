package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/italolelis/coffee-shop/internal/app/order"
	"github.com/italolelis/coffee-shop/internal/pkg/log"
	"go.opencensus.io/trace"
)

type orderHandler struct {
	srv order.Service
}

func (h *orderHandler) router() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.CreateOrder)
	r.Get("/{id}", h.GetOrder)

	return r
}

// CreateCoffee is the handler for order creation
func (h *orderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq struct {
		Name  string      `json:"name"`
		Items order.Items `json:"items"`
	}

	ctx := r.Context()
	logger := log.WithContext(ctx)
	ctx, span := trace.StartSpan(ctx, "handler:create.order", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&orderReq); err != nil {
		http.Error(w, "could not parse the request body", http.StatusBadRequest)
		return
	}

	logger.Debugw("creating a new coffee")
	id, err := h.srv.CreateOrder(ctx, orderReq.Name, orderReq.Items)
	if err != nil {
		if err == order.ErrUnsupportedCoffeeType {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("/orders/%s", id.String()))
	w.WriteHeader(http.StatusCreated)
}

// GetCoffee is the handler for getting a single order
func (h *orderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.WithContext(ctx)
	ctx, span := trace.StartSpan(ctx, "handler:get.order", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	id := chi.URLParam(r, "id")

	data, err := h.srv.RequestOrder(ctx, id)
	if err != nil {
		logger.Errorw("could not get the order", "order_id", id, "err", err)
		if err == order.ErrOrderNotFound {
			http.Error(w, "invalid order id", http.StatusNotFound)
			return
		}

		http.Error(w, "unknown error when retrieving the order", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, data)
}
