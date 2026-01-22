package categories

import (
	"encoding/json"
	"net/http"
	"shadowbane/go-tugas-01/app/helpers"
	"shadowbane/go-tugas-01/app/models"

	"github.com/julienschmidt/httprouter"
)

func Update() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		categoryID := p.ByName("category")

		var updatedCategory models.Category
		err := json.NewDecoder(r.Body).Decode(&updatedCategory)
		if err != nil {
			helpers.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		for i, category := range categoryData {
			if category.ID == categoryID {
				updatedCategory.ID = categoryID
				categoryData[i] = updatedCategory
				helpers.WriteResponse(w, updatedCategory)
				return
			}
		}

		helpers.WriteErrorResponse(w, http.StatusNotFound, "Category not found")
	}
}
