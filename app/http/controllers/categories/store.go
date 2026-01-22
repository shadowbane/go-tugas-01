package categories

import (
	"encoding/json"
	"net/http"
	"shadowbane/go-tugas-01/app/helpers"
	"shadowbane/go-tugas-01/app/models"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func Store() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var newCategory models.Category
		err := json.NewDecoder(r.Body).Decode(&newCategory)
		if err != nil {
			helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		newCategory.ID = strings.ToUpper(helpers.NewULID())

		categoryData = append(categoryData, newCategory)

		helpers.WriteResponse(w, newCategory)
	}
}
