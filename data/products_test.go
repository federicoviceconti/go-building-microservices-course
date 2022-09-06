package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:        "Test",
		Description: "Test",
		Price:       0.1,
		Sku:         "abc123",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
