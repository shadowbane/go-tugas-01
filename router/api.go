package router

import (
	"github.com/julienschmidt/httprouter"

	"github.com/shadowbane/go-tugas-01/handlers"
)

func Api(handlers *handlers.Handler) *httprouter.Router {
	mux := httprouter.New()

	// **GET** `/api/categories` → Get all category
	mux.GET("/api/categories", handlers.CategoryHandler.Index())
	// **POST** `/api/categories` → Add new category
	mux.POST("/api/categories", handlers.CategoryHandler.Store())
	// **PUT** `/api/categories/{id}` → Update category
	mux.PUT("/api/categories/:category", handlers.CategoryHandler.Update())
	// **GET** `/api/categories/{id}` → Get one category
	mux.GET("/api/categories/:category", handlers.CategoryHandler.Show())
	// **DELETE** `/api/categories/{id}` → Delete category
	mux.DELETE("/api/categories/:category", handlers.CategoryHandler.Destroy())

	// **GET** `/api/products` → Get all product
	mux.GET("/api/products", handlers.ProductHandler.Index())
	// **POST** `/api/products` → Add new product
	mux.POST("/api/products", handlers.ProductHandler.Store())
	// **PUT** `/api/products/{id}` → Update product
	mux.PUT("/api/products/:product", handlers.ProductHandler.Update())
	// **GET** `/api/products/{id}` → Get one product
	mux.GET("/api/products/:product", handlers.ProductHandler.Show())
	// **DELETE** `/api/products/{id}` → Delete product
	mux.DELETE("/api/products/:product", handlers.ProductHandler.Destroy())

	// **POST** `/api/checkout` → Checkout / create transaction
	mux.POST("/api/checkout", handlers.TransactionHandler.Checkout())

	return mux
}
