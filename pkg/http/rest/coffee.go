package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/italolelis/coffee-shop/internal/log"
	"github.com/italolelis/coffee-shop/pkg/coffees"
	"go.opencensus.io/trace"
)

type coffeeHandler struct {
	srv coffees.Service
}

func (h *coffeeHandler) router() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.CreateCoffee)
	r.Get("/{id}", h.GetCoffee)

	return r
}

// CreateCoffee is the handler for order creation
func (h *coffeeHandler) CreateCoffee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.WithContext(ctx)
	ctx, span := trace.StartSpan(ctx, "handler:create.coffee", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	var request struct {
		Name  string  `json:"name"`
		Price float32 `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "could not parse the request body", http.StatusBadRequest)
		return
	}

	logger.Debugw("creating a new coffee")
	id, err := h.srv.CreateCoffee(ctx, request.Name, request.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("/coffees/%s", id.String()))
	w.WriteHeader(http.StatusCreated)
}

// GetCoffee is the handler for getting a single order
func (h *coffeeHandler) GetCoffee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.WithContext(ctx)
	ctx, span := trace.StartSpan(ctx, "handler:get.coffee", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	id := chi.URLParam(r, "id")
	coffee, err := h.srv.RequestCoffee(ctx, id)
	if err != nil {
		logger.Errorw("could not get the coffee", "coffee_id", id, "err", err)

		if err != nil {
			logger.Errorw("invalid id", "coffee_id", chi.URLParam(r, "id"))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err == coffees.ErrCoffeeNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "unknown error when retrieving the coffee", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, coffee)
}
