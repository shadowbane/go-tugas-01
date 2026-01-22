package router

import (
	"github.com/julienschmidt/httprouter"

	categoryController "shadowbane/go-tugas-01/app/http/controllers/categories"
)

func Api() *httprouter.Router {
	mux := httprouter.New()

	// **GET** `/categories` → Get all category
	mux.GET("/categories", categoryController.Index())
	// **POST** `/categories` → Add new category
	mux.POST("/categories", categoryController.Store())
	// **PUT** `/categories/{id}` → Update category
	mux.PUT("/categories/:category", categoryController.Update())
	// **GET** `/categories/{id}` → Get one category
	mux.GET("/categories/:category", categoryController.Show())
	// **DELETE** `/categories/{id}` → Delete category
	mux.DELETE("/categories/:category", categoryController.Destroy())

	return mux
}
