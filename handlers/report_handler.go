package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/shadowbane/go-tugas-01/pkg/helpers"
	"github.com/shadowbane/go-tugas-01/repositories"
	"github.com/shadowbane/go-tugas-01/services"
)

type ReportHandler struct {
	service *services.TransactionService
}

func NewReportHandler(service *services.TransactionService) *ReportHandler {
	return &ReportHandler{
		service: service,
	}
}

func (h *ReportHandler) TodayReport() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		includeDetails := r.URL.Query().Get("detail") == "true"

		report, err := h.service.GetTodayReport(includeDetails)
		if err != nil {
			helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.WriteResponse(w, report)
	}
}

func (h *ReportHandler) DateRangeReport() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")

		if startDateStr == "" || endDateStr == "" {
			helpers.WriteErrorResponse(w, http.StatusBadRequest, "start_date and end_date are required")
			return
		}

		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid start_date format, use YYYY-MM-DD")
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			helpers.WriteErrorResponse(w, http.StatusBadRequest, "invalid end_date format, use YYYY-MM-DD")
			return
		}

		// Add one day to endDate to include the entire day (end of day)
		endDate = endDate.AddDate(0, 0, 1)

		includeDetails := r.URL.Query().Get("detail") == "true"
		consolidated := r.URL.Query().Get("consolidated") == "true"

		report, err := h.service.GetReportByDateRange(startDate, endDate, includeDetails, consolidated)
		if err != nil {
			if errors.Is(err, repositories.ErrInvalidDateRange) {
				helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			helpers.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.WriteResponse(w, report)
	}
}
