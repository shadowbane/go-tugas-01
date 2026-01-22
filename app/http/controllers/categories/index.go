package categories

import (
	"io"
	"net/http"
	"shadowbane/go-tugas-01/app/helpers"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func Index() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				zap.S().Errorw("Error while closing body", "error", err)
			}
		}(r.Body)

		helpers.WriteResponse(w, categoryData)
	}
}
