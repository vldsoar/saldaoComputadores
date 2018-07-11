package handlers

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"saldaoComputadores/db"
	"log"
	"saldaoComputadores/helpers"
	"encoding/json"
	"net/rpc"
	"saldaoComputadores/types"
	"strconv"
	"saldaoComputadores/services"
	"regexp"
	"fmt"
)

var (
	gcurrIDSale Counter
	ghostname string
	gclientsHosts []string
)

type Functions interface {
	Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	Cart(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	Purchase(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}

type Control struct {
	View            *template.Template
	ServicesClients []rpc.Client
}

// Instancia um controlador
func Load(host string, lastIDSale int, clientsHosts []string) *Control {
	gcurrIDSale = Counter(lastIDSale + 1)
	ghostname = host
	gclientsHosts = clientsHosts

	c := new(Control)
	c.ServicesClients = services.LoadClientServices(clientsHosts)

	var err error

	c.View, err = template.New("Main").Funcs(helpers.TmplFunctions()).ParseGlob("templates/*.gohtml")

	if err != nil {
		log.Fatalln(err)
	}

	return c
}

// Função para rota de página principal
func (c *Control) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	data, err := db.GetProducts()

	if err != nil {
		log.Fatalln(err)
	}

	c.View.ExecuteTemplate(w, "index.gohtml", data)
}

func (c *Control) Cart(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {
	c.View.ExecuteTemplate(w, "cart.gohtml", nil)
}
// Função para rota de Compra
func (c *Control) Purchase() httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Pega valores da requisição
		r.ParseForm()
		var productsSale types.Products

		for _, v := range ParseFormCollection(r, "products") {
			id, _ := strconv.Atoi(v["id"])
			qtd, _ := strconv.Atoi(v["quantity"])

			product := types.Product{
				ID:id,
				Quantity:qtd,
			}
			// Checa se há estoque local
			if qtd := db.GetProduct(product.ID).Quantity; qtd < product.Quantity {
				str := fmt.Sprintf("Pruduto %d sem estoque", product.ID)
				json.NewEncoder(w).Encode(str)
				return
			}

			productsSale = append(productsSale, product)
		}

		totalPrice, _ := strconv.ParseFloat(r.FormValue("totalPrice"), 32)

		// cria uma nova instancia de Sale
		currSale := types.NewSale(int(gcurrIDSale.GetAndIncrement()))
		currSale.Products = productsSale
		currSale.TotalPrice = float32(totalPrice)

		// Cria uma instancia de Requisição para chamadas remotas
		req := types.Request{
			From:ghostname,
			Clock:services.LClock.Increment(),
			Body:make(map[string]interface{}),
		}

		req.Body["sale"] = currSale

		or := services.Order{}
		resLocal := types.NewResponse()
		or.Check(req, &resLocal)

		if !resLocal.Success {
			json.NewEncoder(w).Encode("Produto(s) sem estoque no momento")
			return
		}

		// Carrega serviços
		c.ServicesClients = services.LoadClientServices(gclientsHosts)
		log.Println("Hosts, ", gclientsHosts)
		log.Println("Clients Services (Purchase): ", c.ServicesClients)

		// Número de respostas esperas
		expectedResponses := len(c.ServicesClients)

		responses := map[int]bool{}
		// Para cada serviço faz uma requisição
		for index, service := range c.ServicesClients {
			res := types.NewResponse()
			err := service.Call("OrderService.Check", req, &res)

			log.Println("Call: ", index, service)

			if err != nil {
				log.Println(err)
			}

			if res.Success {
				responses[index] = true
			}
		}

		// Se o numero de respostas foram iguais as esperadas
		// Registra venda localmente e em todos os serviços conectados
		if len(responses) == expectedResponses {
			db.CreateSale(currSale)
			log.Println("DB")

			for _, service := range c.ServicesClients {
				log.Println("Register Sale Services")
				res := types.NewResponse()
				defer service.Close()
				service.Go("OrderService.RegisterSale", req, &res, nil)
			}

		} else {

			json.NewEncoder(w).Encode("Produto(s) sem estoque no momento")
			return
		}

		json.NewEncoder(w).Encode("Compra efetuada com Sucesso")
	}

}

// Converte requisição em map
func ParseFormCollection(r *http.Request, typeName string) []map[string]string {
	var result []map[string]string
	r.ParseForm()
	for key, values := range r.Form {
		re := regexp.MustCompile(typeName + "\\[([0-9]+)\\]\\[([a-zA-Z]+)\\]")
		matches := re.FindStringSubmatch(key)

		if len(matches) >= 3 {

			index, _ := strconv.Atoi(matches[1])

			for ; index >= len(result); {
				result = append(result, map[string]string{})
			}

			result[index][matches[2]] = values[0]
		}
	}
	return result
}