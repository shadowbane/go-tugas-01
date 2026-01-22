package router

import (
	"github.com/julienschmidt/httprouter"
)

func Api() *httprouter.Router {
	mux := httprouter.New()

	// **GET** `/categories` → Get all category
	mux.GET("/api/v1/categories", nil)
	// **POST** `/categories` → Add new category
	mux.POST("/api/v1/categories", nil)
	// **PUT** `/categories/{id}` → Update category
	mux.PUT("/api/v1/categories/:category", nil)
	// **GET** `/categories/{id}` → Get one category
	mux.GET("/api/v1/categories/:category", nil)
	// **DELETE** `/categories/{id}` → Delete category
	mux.DELETE("/api/v1/categories/:category", nil)

	return mux
}
