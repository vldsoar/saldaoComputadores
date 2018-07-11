package routes

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"saldaoComputadores/handlers"
)

func Load(c *handlers.Control) *httprouter.Router {
	router := httprouter.New()

	router.ServeFiles("/static/*filepath", http.Dir("public/static"))
	router.ServeFiles("/upload/*filepath", http.Dir("public/upload"))

	router.GET("/", c.Index)
	router.GET("/cart", c.Cart)
	router.POST("/purchase", c.Purchase())

	return router
}
