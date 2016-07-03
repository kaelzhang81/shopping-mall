package main

import (
	"infrastructure"
	"interfaces"
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
}
