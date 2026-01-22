package categories

import (
	"github.com/shadowbane/go-tugas-01/app/helpers"
	"io"
	"net/http"

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

		var deleted bool
		categoryID := p.ByName("category")

		func() {
			mutex.Lock()
			defer mutex.Unlock()

			for i, category := range categoryData {
				if category.ID == categoryID {
					categoryData = append(categoryData[:i], categoryData[i+1:]...)
					deleted = true
					return
				}
			}
		}()

		if deleted {
			helpers.WriteResponse(w, map[string]string{"message": "Category deleted"})
		} else {
			helpers.WriteErrorResponse(w, http.StatusNotFound, "Category not found")
		}
	}
}
