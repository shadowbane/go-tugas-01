package categories

import (
	"encoding/json"
	"io"
	"net/http"
	"shadowbane/go-tugas-01/app/helpers"
	"shadowbane/go-tugas-01/app/models"
	"strings"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func Store() httprouter.Handle {
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
		newCategory.ID = strings.ToUpper(helpers.NewULID())

		func() {
			mutex.Lock()
			defer mutex.Unlock()

			categoryData = append(categoryData, newCategory)
		}()

		helpers.WriteResponse(w, newCategory)
	}
}
