package services

import (
	"saldaoComputadores/db"
	"saldaoComputadores/types"
	"log"
)

type Order struct {
	ServiceRemote
}

func (o *Order) Check(req types.Request, res *types.Response) error {

	log.Println(req)

	// Adiciona na fila de eventos
	Events.Put(req)

	events, err := Events.Get(Events.Len())
	qpe := make(map[int]int)
	res.From = o.Host

	if err != nil {
		res.Success = false
		return err
	}

	// acumula a quantidade de produtos, nos eventos
	for _, e := range events {
		r := e.(types.Request)
		s := r.Body["sale"].(types.Sale)

		if r.Clock <= req.Clock {
			for _, p := range s.Products {
				qpe[p.ID] += p.Quantity
			}
		}

		if r.From != req.From && r.Clock != req.Clock {
			Events.Put(e)
		}

	}

	//
	currSale := req.Body["sale"].(types.Sale)

	// Checa se há quantidade de produtos disponíveis
	// considerando as intenções de compra e a quantidade em estoque
	for _, p := range currSale.Products {
		product := db.GetProduct(p.ID)
		if product.Quantity < ( qpe[p.ID] + p.Quantity ) {
			res.Success = false
			return nil
		}
	}

	// Checa e atualiza relógio lógico
	LClock.Witness(req.Clock)

	res.Success = true

	return nil
}

// Registrar Venda
func (o *Order) RegisterSale(req types.Request, res *types.Response) error {

	sale := req.Body["sale"].(types.Sale)

	db.CreateSale(sale)

	res.Success = true
	res.From = o.Host

	return nil
}
