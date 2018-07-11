package db

import (
	js "github.com/thedevsaddam/gojsonq"
	"saldaoComputadores/types"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"os"
	"io"
	"sync"
	"bytes"
)

//var globalStore dbStore
var lock sync.Mutex

// Carrega um json com o arquivo especificado
func Load(name string) *js.JSONQ {
	db := js.New().File(fmt.Sprintf("db/%s.json", name))

	if db.Error() != nil {
		return nil
	}

	return db
}

// Recupera todos os produtos
func GetProducts() (types.Products, error) {

	db := Load("products")

	bn, _  := json.Marshal(db.From("products").Get())

	products := types.Products{}

	json.Unmarshal(bn, &products)

	return products, nil
}

// Recupera um produto por seu id
func GetProduct(id int) types.Product {
	lock.Lock()
	defer lock.Unlock()
	db := Load("products")

	findProduct := db.From("products").WhereEqual("id", id).First()

	bn, _  := json.Marshal(findProduct)

	product := types.Product{}

	json.Unmarshal(bn, &product)

	return product
}

// Altera todos os produtos no arquivo
func SetProducts(v interface{})  {

	f, _ := ioutil.ReadFile("db/products.json")

	pjson, _ := simplejson.NewJson(f)

	pjson.Set("products", v)

	Save("db/products.json", pjson)

}

// Cria uma compra e salva no arquivo
func CreateSale(v interface{})  {

	f, _ := os.Open("db/sales.json")

	pjson, _ := simplejson.NewFromReader(f)

	arr, _ := pjson.Array()

	sales := append(arr, v)

	Save("db/sales.json", sales)

	productsDB, _ := GetProducts()

	sale := v.(types.Sale)

	for _, iProduct := range sale.Products {
		for i, product := range productsDB {
			if iProduct.ID == product.ID {
				(productsDB)[i].Quantity -= iProduct.Quantity
				break
			}
		}
	}

	SetProducts(productsDB)

}

// Save saves a representation of v to the file at path.
func Save(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := Marshal(v)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}

func Marshal(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}