package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
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
	case http.MethodGet:
		GetProducts(writer)
	case http.MethodPut:
		UpdateProduct(writer, request)
	case http.MethodPost:
		AddProduct(writer, request)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}

	return
}

func AddProduct(writer http.ResponseWriter, request *http.Request) {
	prod := &data.Product{}
	err := prod.FromJson(request.Body)

	if err != nil {
		http.Error(writer, "Product body is not well-formed", http.StatusBadRequest)
		return
	} else {
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
}

func UpdateProduct(writer http.ResponseWriter, request *http.Request) {
	r := regexp.MustCompile("/([0-9]+)")

	matches := r.FindAllStringSubmatch(request.URL.Path, -1)

	if len(matches) > 0 {
		id, errParse := strconv.Atoi(matches[0][1])

		if errParse != nil {
			http.Error(writer, "Exception on retrieving id", http.StatusBadRequest)
			return
		}

		prod := &data.Product{}
		err := prod.FromJson(request.Body)

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
}

func GetProducts(writer http.ResponseWriter) {
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

// Encode the object into JSON object. See the Product object to know how to build it
//jsonValue, err := json.Marshal(lp)
