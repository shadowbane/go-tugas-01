package helpers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type ResponseData struct {
	Success bool         `json:"success"`
	Data    *interface{} `json:"data"`
}

type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func (r *ResponseData) SetData(data interface{}) {
	r.Data = &data
}

func (r *ResponseData) ToJson() []byte {
	response, _ := json.Marshal(r)
	return response
}

func WriteResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	responseData := ResponseData{
		Success: true,
		Data:    nil,
	}

	responseData.SetData(data)

	_, err := w.Write(responseData.ToJson())
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

func WritePaginatedResponse(w http.ResponseWriter, items interface{}, pagination Pagination) {
	paginatedData := PaginatedResponse{
		Items:      items,
		Pagination: pagination,
	}
	WriteResponse(w, paginatedData)
}

// WriteErrorResponse writes an error response to the client in the following format:
//
//	{
//	  "success": false,
//	  "data": {
//	    "message": "ERRINFO"
//	  }
//	}
func WriteErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)

	responseData := &ResponseData{
		Success: false,
	}

	responseData.SetData(map[string]string{
		"message": errorMsg,
	})

	err := json.NewEncoder(w).Encode(responseData)

	if err != nil {
		zap.S().Fatalf(err.Error())
	}
}
