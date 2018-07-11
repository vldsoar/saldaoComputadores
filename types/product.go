package types

import "reflect"

type Product struct {
	ID			int 	`json:"id"`
	Name 		string 	`json:"name,omitempty"`
	Image 		string 	`json:"image,omitempty"`
	Price 		float32 `json:"price,omitempty"`
	Quantity 	int 	`json:"quantity"`
}

type Products []Product

func (p Product) IsEmpty() bool {
	return reflect.DeepEqual(p, Product{})
}