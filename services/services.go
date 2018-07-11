package services

import (
	"saldaoComputadores/types"
	"log"
	"github.com/golang-collections/go-datastructures/queue"
	"net/rpc"
)

var LClock *types.LamportClock
var Events *queue.PriorityQueue

func init()  {
	LClock = &types.LamportClock{}
	Events = queue.NewPriorityQueue(10)
}

type Service interface {
	Ping(req types.Request, res *types.Response) error
}

type ServiceRemote struct {
	Host string
}

func (s *ServiceRemote) Ping(req types.Request, res *types.Response) error {
	res.Success = true
	res.From = s.Host
	log.Println("Ping, ", req.From)
	return nil
}


func LoadClientServices(hosts []string) []rpc.Client {
	log.Println("Loading Clients Services")

	var cServices []rpc.Client

	log.Println(hosts)

	for ind, host := range hosts {
		log.Println(ind)
		//go func(h string) {
			client, err := rpc.DialHTTPPath("tcp", host, "/rpc")

			if err == nil {
				cServices = append(cServices, *client)
			} else {
				log.Print("LoadClientServices: ")
				log.Println(err)
			}
		//}(host)

	}

	return cServices
}
