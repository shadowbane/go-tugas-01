package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/shadowbane/go-tugas-01/models"
	"github.com/shadowbane/go-tugas-01/pkg/helpers"
	"github.com/shadowbane/go-tugas-01/repositories"
	"github.com/shadowbane/go-tugas-01/services"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) Index() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				zap.S().Errorw("Error while closing body", "error", err)
			}
		}(r.Body)

		name := r.URL.Query().Get("name")

		data, err := h.service.GetAll(name)
		if err != nil {
			helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.WriteResponse(w, data)
	}
}

func (h *ProductHandler) Show() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				zap.S().Errorw("Error while closing body", "error", err)
			}
		}(r.Body)

		productID := p.ByName("product")
		product, err := h.service.GetByID(productID)
		if err != nil {
			if errors.Is(err, repositories.ErrProductNotFound) {
				helpers.WriteErrorResponse(w, http.StatusNotFound, "Product not found")
				return
			}
			helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.WriteResponse(w, product)
	}
}

func (h *ProductHandler) Store() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				zap.S().Errorw("Error while closing body", "error", err)
			}
		}(r.Body)

		var newProduct models.Product
		err := json.NewDecoder(r.Body).Decode(&newProduct)
		if err != nil {
			helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = h.service.Create(&newProduct)
		if err != nil {
			helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.WriteResponse(w, newProduct)
	}
}

func (h *ProductHandler) Update() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				zap.S().Errorw("Error while closing body", "error", err)
			}
		}(r.Body)

		productID := p.ByName("product")

		var updatedProduct models.Product
		err := json.NewDecoder(r.Body).Decode(&updatedProduct)
		if err != nil {
			helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = h.service.Update(productID, &updatedProduct)
		if err != nil {
			if errors.Is(err, repositories.ErrProductNotFound) {
				helpers.WriteErrorResponse(w, http.StatusNotFound, "Product not found")
				return
			}
			helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.WriteResponse(w, updatedProduct)
	}
}

func (h *ProductHandler) Destroy() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				zap.S().Errorw("Error while closing body", "error", err)
			}
		}(r.Body)

		productID := p.ByName("product")
		err := h.service.Delete(productID)
		if err != nil {
			if errors.Is(err, repositories.ErrProductNotFound) {
				helpers.WriteErrorResponse(w, http.StatusNotFound, "Product not found")
				return
			}
			helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.WriteResponse(w, map[string]string{"message": "Product deleted"})
	}
}
