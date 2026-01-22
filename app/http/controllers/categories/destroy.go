package categories

import (
	"io"
	"net/http"
	"shadowbane/go-tugas-01/app/helpers"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func Destroy() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				zap.S().Errorw("Error while closing body", "error", err)
			}
		}(r.Body)

		categoryID := p.ByName("category")

		for i, category := range categoryData {
			if category.ID == categoryID {
				categoryData = append(categoryData[:i], categoryData[i+1:]...)
				helpers.WriteResponse(w, map[string]string{"message": "Category deleted"})
				return
			}
		}

		helpers.WriteErrorResponse(w, http.StatusNotFound, "Category not found")
	}
}
