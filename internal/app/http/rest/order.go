package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/italolelis/coffee-shop/internal/app/order"
	"github.com/italolelis/coffee-shop/internal/pkg/log"
)

type OrderHandler struct {
	srv order.Service
}

func (h OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		logger = log.WithContext(ctx).Named("orders").With("action", "checkout")
	)

	var cmd order.CheckoutCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		logger.Errorw("failed to decode payload", "err", err)

		http.Error(w, "failed to decode payload", http.StatusBadRequest)

		return
	}

	orderID, err := h.srv.Checkout(ctx, cmd)
	if err != nil {
		logger.Errorw("failed to checkout order", "err", err)

		http.Error(w, "failed to checkout order", http.StatusInternalServerError)

		return
	}

	w.Header().Add("Location", fmt.Sprintf("/orders/%s", orderID))
	w.WriteHeader(http.StatusCreated)
}

func (h OrderHandler) AddToOrder(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		logger = log.WithContext(ctx).Named("orders").With("action", "add-to-order")
	)

	var cmd order.AddToOrderCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		logger.Errorw("failed to decode payload", "err", err)

		http.Error(w, "failed to decode payload", http.StatusBadRequest)

		return
	}

	orderID, err := h.srv.AddToOrder(ctx, cmd)
	if err != nil {
		logger.Errorw("failed to add items to order", "err", err)

		http.Error(w, "failed to add items to order", http.StatusInternalServerError)

		return
	}

	w.Header().Add("Location", fmt.Sprintf("/orders/%s", orderID))
	w.WriteHeader(http.StatusCreated)
}

func (h OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	var (
		ctx        = r.Context()
		logger     = log.WithContext(ctx).Named("orders").With("action", "get-order")
		rawOrderID = chi.URLParam(r, "orderID")
	)

	orderID, err := uuid.Parse(rawOrderID)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	o, err := h.srv.Fetch(ctx, orderID)
	if err != nil {
		if errors.Is(order.ErrNotFound, err) {
			http.Error(w, "couldn't find order", http.StatusNotFound)
			return
		}

		logger.Errorw("failed to fetch order", "err", err)
		http.Error(w, "failed to fetch order", http.StatusInternalServerError)

		return
	}

	render.JSON(w, r, o)
}
