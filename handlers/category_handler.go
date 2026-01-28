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

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

func (h *CategoryHandler) Index() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				zap.S().Errorw("Error while closing body", "error", err)
			}
		}(r.Body)

		data, err := h.service.GetAll()
		if err != nil {
			helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.WriteResponse(w, data)
	}
}

func (h *CategoryHandler) Show() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				zap.S().Errorw("Error while closing body", "error", err)
			}
		}(r.Body)

		categoryID := p.ByName("category")
		category, err := h.service.GetByID(categoryID)

		if err != nil {
			if errors.Is(err, repositories.ErrCategoryNotFound) {
				helpers.WriteErrorResponse(w, http.StatusNotFound, "Category not found")
				return
			}
			helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.WriteResponse(w, category)
	}
}

func (h *CategoryHandler) Store() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				zap.S().Errorw("Error while closing body", "error", err)
			}
		}(r.Body)

		var newCategory models.Category
		err := json.NewDecoder(r.Body).Decode(&newCategory)
		if err != nil {
			helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = h.service.Create(&newCategory)
		if err != nil {
			helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.WriteResponse(w, newCategory)
	}
}

func (h *CategoryHandler) Update() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				zap.S().Errorw("Error while closing body", "error", err)
			}
		}(r.Body)

		var updatedCategory models.Category
		err := json.NewDecoder(r.Body).Decode(&updatedCategory)
		if err != nil {
			helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		categoryID := p.ByName("category")
		err = h.service.Update(categoryID, &updatedCategory)

		if err != nil {
			if errors.Is(err, repositories.ErrCategoryNotFound) {
				helpers.WriteErrorResponse(w, http.StatusNotFound, "Category not found")
				return
			}
			helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.WriteResponse(w, updatedCategory)
	}
}

func (h *CategoryHandler) Destroy() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				zap.S().Errorw("Error while closing body", "error", err)
			}
		}(r.Body)

		categoryID := p.ByName("category")
		err := h.service.Delete(categoryID)

		if err != nil {
			if errors.Is(err, repositories.ErrCategoryNotFound) {
				helpers.WriteErrorResponse(w, http.StatusNotFound, "Category not found")
				return
			}
			helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.WriteResponse(w, map[string]string{"message": "Category deleted"})
	}
}
