package main

import (
	"infrastructure"
	"interfaces"
	"net/http"
	"usecases"
)

func main() {
	dbHandler := infrastructure.NewSqliteHandler("/var/tmp/production.sqlite")

	handlers := make(map[string]interfaces.DbHandler)
	handlers["DbUserRepo"] = dbHandler
	handlers["DbCustomerRepo"] = dbHandler
	handlers["DbItemRepo"] = dbHandler
	handlers["DbOrderRepo"] = dbHandler

	orderInteractor := new(usecases.OrderInteractor)
	orderInteractor.UserRepository = interfaces.NewDbUserRepo(handlers)
	orderInteractor.OrderRepository = interfaces.NewDbOrderRepo(handlers)
	orderInteractor.ItemRepository = interfaces.NewDbItemRepo(handlers)

	webserviceHandler := interfaces.WebServiceHandler{}
	webserviceHandler.OrderInteractor = orderInteractor

	http.HandleFunc("/orders", func(res http.ResponseWriter, req *http.Request) {
		webserviceHandler.ShowOrder(res, req)
	})

	http.ListenAndServe(":8080", nil)
}
