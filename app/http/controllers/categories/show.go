package categories

import (
	"github.com/shadowbane/go-tugas-01/app/helpers"
	"github.com/shadowbane/go-tugas-01/app/models"
	"io"
	"net/http"

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
		var foundCategory models.Category
		var found bool

		func() {
			mutex.RLock()
			defer mutex.RUnlock()

			for _, category := range categoryData {
				if category.ID == categoryID {
					foundCategory = category
					found = true
					return
				}
			}
		}()

		if found {
			helpers.WriteResponse(w, foundCategory)
		} else {
			helpers.WriteErrorResponse(w, http.StatusNotFound, "Category not found")
		}
	}
}
