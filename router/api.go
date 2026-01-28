package router

import (
	"github.com/julienschmidt/httprouter"

	"github.com/shadowbane/go-tugas-01/handlers"
)

func Api(handlers *handlers.Handler) *httprouter.Router {
	mux := httprouter.New()

	// **GET** `/categories` → Get all category
	mux.GET("/categories", handlers.CategoryHandler.Index())
	// **POST** `/categories` → Add new category
	mux.POST("/categories", handlers.CategoryHandler.Store())
	// **PUT** `/categories/{id}` → Update category
	mux.PUT("/categories/:category", handlers.CategoryHandler.Update())
	// **GET** `/categories/{id}` → Get one category
	mux.GET("/categories/:category", handlers.CategoryHandler.Show())
	// **DELETE** `/categories/{id}` → Delete category
	mux.DELETE("/categories/:category", handlers.CategoryHandler.Destroy())

	// **GET** `/products` → Get all product
	mux.GET("/products", handlers.ProductHandler.Index())
	// **POST** `/products` → Add new product
	mux.POST("/products", handlers.ProductHandler.Store())
	// **PUT** `/products/{id}` → Update product
	mux.PUT("/products/:product", handlers.ProductHandler.Update())
	// **GET** `/products/{id}` → Get one product
	mux.GET("/products/:product", handlers.ProductHandler.Show())
	// **DELETE** `/products/{id}` → Delete product
	mux.DELETE("/products/:product", handlers.ProductHandler.Destroy())

	return mux
}
