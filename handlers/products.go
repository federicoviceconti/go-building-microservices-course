package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"services/data"
	"strconv"
	"time"
)

type ProductHandler struct {
	l *log.Logger
}

func (p *ProductHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	println("Handling", request.Method, "on path -> \"/\"")

	switch request.Method {
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func (p *ProductHandler) AddProduct(writer http.ResponseWriter, request *http.Request) {
	prod := request.Context().Value(KeyProduct{}).(*data.Product)
	prod.CreatedOn = time.Now().String()
	prod.UpdatedOn = time.Now().String()

	fmt.Printf("Product to add: %#v", prod)

	data.AddProduct(prod)

	products := make(data.Products, 1)
	products = append(products, prod)

	writer.WriteHeader(200)
	err := products.ToJson(writer)
	if err != nil {
		http.Error(writer, "Adding a product makes an exception", http.StatusBadRequest)
		return
	}
}

func (p *ProductHandler) UpdateProduct(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		prod := request.Context().Value(KeyProduct{}).(*data.Product)

		if err != nil {
			http.Error(writer, "Product body is not well-formed", http.StatusBadRequest)
			return
		}

		errProduct := data.UpdateProduct(prod, id)
		if errProduct != nil {
			message := fmt.Sprintf("Product with id %d not found", id)

			http.Error(writer, message, http.StatusNotFound)
			return
		}
	} else {
		http.Error(writer, "\"id\" field is mandatory on URI. E.g. /{id}", http.StatusBadRequest)
	}

	writer.WriteHeader(200)
}

func (p *ProductHandler) GetProducts(writer http.ResponseWriter, request *http.Request) {
	lp := data.GetProducts()

	// We can use an encoder and avoid to use Marshal
	err := lp.ToJson(writer)

	if err != nil {
		http.Error(writer, "An error occur", http.StatusBadRequest)
	}
}

func NewProductHandler(l *log.Logger) *ProductHandler {
	return &ProductHandler{l}
}

type KeyProduct struct{}

func (p *ProductHandler) MiddlewareProductsValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		prod := &data.Product{}
		err := prod.FromJson(request.Body)

		if err != nil {
			http.Error(writer, "Product body is not well-formed", http.StatusBadRequest)
			return
		}

		validateError := prod.Validate()
		if validateError != nil {
			errorMessage := fmt.Sprintf("Product values are not valid %s", validateError)
			http.Error(writer, errorMessage, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(request.Context(), KeyProduct{}, prod)
		req := request.WithContext(ctx)

		next.ServeHTTP(writer, req)
	})
}
