package main

import (
	"saldaoComputadores/types"
	"log"
	"net/rpc"
	"time"
	"math/rand"
	"fmt"
)

func main()  {

	//gob.Register(types.Sale{})
	//gob.Register(types.Sales{})
	//gob.Register(types.Products{})
	//gob.Register(types.Product{})
	//
	//url := "192.168.0.104:8080"
	//
	//host := flag.String("host", "default", "default")
	//flag.Parse()
	//
	////client, err := rpc.DialHTTP("tcp", url)
	//
	//client, err := rpc.DialHTTPPath("tcp", url, "/rpc")
	//
	//
	//if err != nil {
	//	log.Fatal("dialing:", err)
	//}
	//
	//res := types.Response{}
	//
	//body := make(map[string]interface{})
	//
	//pds := types.Products{
	//	types.Product{
	//		ID:1,
	//		Quantity:1,
	//	},
	//}
	//
	//log.Println(pds)
	//
	//sale := types.Sale{}
	//
	//sale = types.Sale{
	//	ID:1,
	//	Host:*host,
	//	Datetime:time.Now(),
	//	Products:pds,
	//	TotalPrice:1259.00,
	//}
	//
	//log.Println(sale)
	//
	//body["sale"] = sale
	//body["lastSale"] = 1
	//
	//req := types.Request{
	//	From: *host,
	//	Body: body,
	//}
	//
	//
	//for {
	//	select {
	//	case <-time.After(time.Second * 5):
	//		err = client.Call("SyncDBService.Ping", req, &res )
	//
	//		if err != nil {
	//			log.Fatal("SyncDBService error:", err)
	//		}
	//
	//		log.Println(res)
	//	}
	//}


	lclock := &types.LamportClock{}


	//
	//req1 := &types.Request{}
	//
	//req1.Clock = lclock.Time()
	//
	//lclock.Increment()
	//
	//req2 := &types.Request{}
	//req2.Clock = lclock.Time()
	//
	//q := queue.NewPriorityQueue(10)
	//
	//q.Put(req1, req2)
	//
	//
	//log.Println(q.Peek())
	//
	//log.Println(q.Peek())

	url := "192.168.5.107:8085"

	client, err := rpc.DialHTTPPath("tcp", url, "/rpc")

	if err != nil {
		log.Fatal("dialing:", err)
	}

	res := types.Response{}
	req := types.NewRequest()

	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(100)

	req.From = fmt.Sprintf("%s - %d", "test ", i)

	for {
		select {
		case <-time.After(time.Second * 5):
			lclock.Increment()
			req.Clock = lclock.Time()

			err := client.Call("OrderService.Ping", req, &res)

			if err != nil {
				panic(err)
			}

			log.Println(req)
			log.Println(res.Success)
		}
	}
}
