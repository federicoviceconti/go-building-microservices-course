// Package handlers Products API.
//
// Documentation for Product API
//
//	    Schemes: http
//	    BasePath: /
//	    Version: 1.0.0
//
//	    Consumes:
//		- application/json
//
//	    Produces:
//		- application / json
//
// swagger:meta
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

// productsResponseWrapper list of products, which returns into response.
// swagger:response productsResponse
type productsResponseWrapper struct {
	// The product list
	// in: body
	Body []data.Product
}

// swagger:parameters deleteProduct
type productDeleteParametersWrapper struct {
	// The product id
	// required: true
	// in: path
	id int
}

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
	prod := request.Context().Value(keyProduct{}).(*data.Product)
	prod.CreatedOn = time.Now().String()
	prod.UpdatedOn = time.Now().String()

	fmt.Printf("Product to add: %#v", prod)

	data.AddProduct(prod)

	products := make(data.Products, 1)
	products = append(products, prod)

	err := products.ToJson(writer)
	if err != nil {
		http.Error(writer, "Adding a product makes an exception", http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusCreated)
}

// swagger:route DELETE /products/{id} listProducts
//
// Delete a product, with the given id.
// Returns list of products without the deleted one.
//
//		Consumes:
//	    - application/json
//
//		Produces:
//	    - application/json
//
//		Schemes: http
//
//		Parameters:
//			+	name:id
//				in: query
//				type: integer
//				required: true
//
//		Responses:
//			200: productsResponse
//
// DeleteProduct deletes a product from the data source
func (p *ProductHandler) DeleteProduct(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id, errParse := strconv.Atoi(vars["id"])

	if errParse != nil {
		if !data.DeleteProductById(id) {
			http.Error(writer, "products does not exists", http.StatusBadRequest)
		}
	} else {
		http.Error(writer, "\"id\" field is mandatory on URI. E.g. /{id}", http.StatusBadRequest)
	}

	writer.WriteHeader(http.StatusOK)
	errGet := data.GetProducts().ToJson(writer)
	if errGet != nil {
		http.Error(writer, "An exception occur on parsing data", http.StatusBadRequest)
	}
}

func (p *ProductHandler) UpdateProduct(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		prod := request.Context().Value(keyProduct{}).(*data.Product)

		errProduct := data.UpdateProduct(prod, id)
		if errProduct != nil {
			message := fmt.Sprintf("Product with id %d not found", id)

			http.Error(writer, message, http.StatusNotFound)
			return
		}
	} else {
		http.Error(writer, "\"id\" field is mandatory on URI. E.g. /{id}", http.StatusBadRequest)
	}

	writer.WriteHeader(http.StatusOK)
}

// swagger:route GET /products listProducts
//
// Returns a list of products
//
//	Responses:
//	  200: productsResponse
//
// GetProducts returns a list of products
func (p *ProductHandler) GetProducts(writer http.ResponseWriter, _ *http.Request) {
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

type keyProduct struct{}

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

		ctx := context.WithValue(request.Context(), keyProduct{}, prod)
		req := request.WithContext(ctx)

		next.ServeHTTP(writer, req)
	})
}
