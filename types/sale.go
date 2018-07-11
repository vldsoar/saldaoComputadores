package types

import (
	"time"
	"reflect"
)

type Sale struct {
	ID 			int 		`json:"id"`
	Host 		string 		`json:"host"`
	Datetime 	time.Time 	`json:"datetime"`
	Products 	Products 	`json:"products"`
	TotalPrice 	float32 	`json:"totalPrice"`
}

type Sales []Sale

func (s Sale) IsEmpty() bool {
	return reflect.DeepEqual(s, Sale{})
}

func NewSale(id int) Sale {
	return Sale{
		ID:id,
		Datetime:time.Now(),
	}
}