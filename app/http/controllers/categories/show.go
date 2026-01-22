package categories

import (
	"net/http"
	"shadowbane/go-tugas-01/app/helpers"

	"github.com/julienschmidt/httprouter"
)

func Show() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
