package coffees

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/italolelis/kit/log"
	"github.com/satori/go.uuid"
	"go.opencensus.io/trace"
)

type Handler struct {
	wRepo WriteRepository
	rRepo ReadRepository
}

// NewHandler creates a new instance of OrderHandler
func NewHandler(wRepo WriteRepository, rRepo ReadRepository) *Handler {
	return &Handler{wRepo: wRepo, rRepo: rRepo}
}

// CreateCoffee is the handler for order creation
func (h *Handler) CreateCoffee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.WithContext(ctx)
	ctx, span := trace.StartSpan(ctx, "handler:create.coffee", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	c := NewCoffee()
	d := json.NewDecoder(r.Body)
	if err := d.Decode(c); err != nil {
		http.Error(w, "could not parse the request body", http.StatusBadRequest)
		return
	}

	logger.Debugw("creating a new coffee", "coffee_id", c.ID)
	if err := h.wRepo.Add(ctx, c); err != nil {
		http.Error(w, "could not save your coffee", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("/coffees/%s", c.ID.String()))
	w.WriteHeader(http.StatusCreated)
}

// GetCoffee is the handler for getting a single order
func (h *Handler) GetCoffee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.WithContext(ctx)
	ctx, span := trace.StartSpan(ctx, "handler:get.coffee", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		logger.Errorw("invalid id", "coffee_id", chi.URLParam(r, "id"))
		http.Error(w, "invalid coffee id provided", http.StatusBadRequest)
		return
	}

	data, err := h.rRepo.FindOneByID(ctx, id)
	if err != nil {
		logger.Errorw("could not get the coffee", "coffee_id", id, "err", err)
		if err == ErrCoffeeNotFound {
			http.Error(w, "invalid order id", http.StatusNotFound)
			return
		}

		http.Error(w, "unknown error when retrieving the coffee", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, data)
}
