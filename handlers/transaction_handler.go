package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/shadowbane/go-tugas-01/dto"
	"github.com/shadowbane/go-tugas-01/pkg/helpers"
	"github.com/shadowbane/go-tugas-01/repositories"
	"github.com/shadowbane/go-tugas-01/services"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

func (h *TransactionHandler) Checkout() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				zap.S().Errorw("Error while closing body", "error", err)
			}
		}(r.Body)

		var request dto.CheckoutRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		transaction, err := h.service.Checkout(&request)
		if err != nil {
			if errors.Is(err, repositories.ErrEmptyCart) || errors.Is(err, repositories.ErrInvalidQuantity) {
				helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			if errors.Is(err, repositories.ErrProductNotFound) {
				helpers.WriteErrorResponse(w, http.StatusNotFound, "Product not found")
				return
			}
			if errors.Is(err, repositories.ErrInsufficientStock) {
				helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.WriteResponse(w, transaction)
	}
}
