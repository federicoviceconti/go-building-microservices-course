package data

import (
	"encoding/json"
	"io"
	"time"
)

type Product struct {
	Id          int     `json:"id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float32 `json:"price,omitempty"`
	Sku         string  `json:"sku,omitempty"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
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
