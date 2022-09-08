package data

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
	"regexp"
	"time"
)

// Product defines the structure for the Product API
// swagger:model
type Product struct {
	Id int `json:"id,omitempty"`
	// required: true
	Name string `json:"name,omitempty" validate:"required"`
	// required: true
	Description string `json:"description,omitempty" validate:"required"`
	// the product price
	// min: 0
	Price     float32 `json:"price,omitempty" validate:"gt=0"`
	Sku       string  `json:"sku,omitempty" validate:"required,sku-validation"`
	CreatedOn string  `json:"-"`
	UpdatedOn string  `json:"-"`
	DeletedOn string  `json:"-"`
}

func (p *Product) Validate() error {
	validate := validator.New()
	err := validate.RegisterValidation("sku-validation", ValidateSKU)
	if err != nil {
		return err
	}
	return validate.Struct(p)
}

func ValidateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile("[a-z][a-z][a-z]+")
	matches := re.FindAllString(fl.Field().String(), -1)

	return !(len(matches) != 1)
}

// Products defining it to make the code readable
type Products []*Product

// ToJson it's an easy way to convert the products slice into JSON
func (p Products) ToJson(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJson(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.Id = GetNextId()
	productList = append(productList, p)
}

func DeleteProductById(id int) bool {
	index, err := FindProductIndex(id)

	if index > -1 && err != nil {
		productList = append(productList[:index], productList[index+1:]...)
		return true
	}

	return false
}

func UpdateProduct(p *Product, id int) error {
	product, err := FindProduct(id)
	if err != nil {
		return err
	}

	edited := false

	if IsNotEmpty(p.Name) {
		product.Name = p.Name
		edited = true
	}

	if IsNotEmpty(p.Sku) {
		product.Sku = p.Sku
		edited = true
	}

	if IsNotEmpty(p.Description) {
		product.Description = p.Description
		edited = true
	}

	if p.Price > 0 {
		product.Price = p.Price
		edited = true
	}

	if edited {
		product.UpdatedOn = time.Now().String()
	}

	return nil
}

type ProductNotFoundError struct {
	message string
	id      int
}

func (p ProductNotFoundError) Error() string {
	return "product not found"
}

func FindProduct(id int) (*Product, error) {
	for _, product := range productList {
		if product.Id == id {
			return product, nil
		}
	}

	return nil, ProductNotFoundError{"product not found", id}
}

func FindProductIndex(id int) (int, error) {
	for index, product := range productList {
		if product.Id == id {
			return index, nil
		}
	}

	return -1, ProductNotFoundError{"product not found", id}
}

func IsNotEmpty(value string) bool {
	return len(value) > 0
}

func GetNextId() int {
	lastProduct := productList[len(productList)-1]
	return lastProduct.Id + 1
}

// Our data source, for testing purpose
var productList = []*Product{
	{
		Id:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		Sku:         "abc232",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		Id:          2,
		Name:        "Espresso",
		Description: "Short coffee without milk",
		Price:       1.99,
		Sku:         "fjc123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
