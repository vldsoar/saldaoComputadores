package services

import (
	"saldaoComputadores/db"
	"saldaoComputadores/types"
	"log"
)

//type SyncDBService interface {
//	Call(req types.Request, res *types.Response) error
//}

type SyncDB struct {
	ServiceRemote
}

func (s *SyncDB) Call(req types.Request, res *types.Response) error {

	log.Println("Request receive SynDB: ", req)

	lastSale := req.Body["lastSale"].(int)

	salesDB := db.Load("sales").Where("id", ">", lastSale).Get()

	log.Println("salesDB: ", salesDB)

	res.Body = make(map[string]interface{})
	res.Body["sales"] = salesDB

	res.Success = true


	return nil
}