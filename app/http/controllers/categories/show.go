package categories

import (
	"io"
	"net/http"
	"shadowbane/go-tugas-01/app/helpers"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func Show() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				zap.S().Errorw("Error while closing body", "error", err)
			}
		}(r.Body)

		categoryID := p.ByName("category")

		for _, category := range categoryData {
			if category.ID == categoryID {
				helpers.WriteResponse(w, category)
				return
			}
		}

		helpers.WriteErrorResponse(w, http.StatusNotFound, "Category not found")
	}
}
