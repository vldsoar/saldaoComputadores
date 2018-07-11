package config

import (
	"saldaoComputadores/services"
	"saldaoComputadores/handlers"
	"github.com/julienschmidt/httprouter"
	"saldaoComputadores/routes"
	"net/rpc"
	"log"
	"encoding/gob"
	"saldaoComputadores/types"
	"saldaoComputadores/db"
)

type Application struct {
	Config   		Configuration         `json:"config"`
	Router   		*httprouter.Router
	Control			*handlers.Control
}

// Registra os tipos que podem ser passados remotamente
func init()  {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	gob.Register(types.Sale{})
	gob.Register(types.Sales{})
	gob.Register(types.Products{})
	gob.Register(types.Product{})
}

func LoadApplication() *Application {
	// Ler arquivo de cofiguração
	conf, err := ReadConfig("config.json")

	if err != nil {
		panic(err)
	}
	// Instancia uma aplicação
	app := new(Application)

	// inicializa Clientes RPC
	cServices := services.LoadClientServices(conf.HostServices)

	log.Println(cServices)

	// Sincroniza banco de dados
	SyncDB(cServices)

	idSale := db.Load("sales").Max("id")

	control := handlers.Load(conf.Hostname, int(idSale), conf.HostServices)

	*app = Application{
		Config:   conf,
		Router:   routes.Load(control),
		Control:  control,
	}

	registerServices(app)

	return app
}
// Função para sincronizar o banco de dados
func SyncDB(clients []rpc.Client) {

	if len(clients) == 0 {
		return
	}

	syncronized := bool(false)

	for _, s := range clients {

		if !syncronized {
			req := types.NewRequest()
			res := types.NewResponse()

			req.Body["lastSale"] = int(db.Load("sales").Max("id"))

			err := s.Call("SyncDBService.Call", req, &res)

			if err != nil {
				log.Println(err)
				continue
			}

			if len(res.Body["sales"].([]interface{})) > 0 {
				sales := res.Body["sales"].(types.Sales)

				for _, sale := range sales {
					db.CreateSale(sale)
				}

				syncronized = true
				log.Println("Synchronized database")
			}



		} else {
			break
		}
	}

	if !syncronized {
		log.Println("Unsynchronized database")
	}
}

// Registra os serviços rpc
func registerServices(app *Application)  {
	rpcServ := rpc.NewServer()
	sr := services.ServiceRemote{Host:app.Config.Hostname}

	syncDB := &services.SyncDB{sr}
	order := &services.Order{sr}

	err := rpcServ.RegisterName("OrderService", order)
	err = rpcServ.RegisterName("SyncDBService", syncDB)

	if err != nil {
		log.Fatal("Format of service isn't correct. ", err)
	}

	app.Router.Handler("CONNECT", "/rpc", rpcServ)

}


