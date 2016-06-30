package interfaces

import (
	"domain"
	"fmt"
	"usecases"
)

type DbHandler interface {
	Execute(statment string)
	Query(statement string)
}

type Row interface {
	Scan(dest ...interface{})
	Next() bool
}

type DbRepo struct {
	dbHandlers map[string]DbHandler
	dbHandler  DbHandler
}

type DbUserRepo DbRepo
type DbCustomerRepo DbRepo
type DbOrderRepo DbRepo
type DbItemRepo DbRepo

func NewDbUserRepo(DbHandlers map[string]DbHandler) *DbUserRepo {
	dbUserRepo := new(DbUserRepo)
	dbUserRepo.dbHandlers = DbHandlers
	dbUserRepo.dbHandler = DbHandlers["DbUserRepo"]
	return dbUserRepo
}

func (repo *DbUserRepo) Store(user usecases.User) {
	IsAdmin := "no"
	if user.IsAdmin {
		IsAdmin = "yes"
	}

	repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO users (id, customer_id, is_admin) \
    VALUES ('%d', '%d', '%v')`, user.Id, user.Customer.Id, IsAdmin))

	dbCustomerRepo := NewDbCustomerRepo(repo.dbHandlers)
	dbCustomerRepo.Store(user.Customer)
}

func NewDbCustomerRepo(dbHandlers map[string]DbHandler) *DbCustomerRepo {
	repo := new(DbCustomerRepo)
	repo.dbHandlers = dbHandlers
	repo.dbHandler = dbHandlers["DbCustomerRepo"]
	return repo
}

func (repo *DbCustomerRepo) Store(customer domain.Customer) {
	repo.dbHandler.Exexute(fmt.Sprintf(`INSERT INTO customers (id, name) 
                                        VALUES ('%d', '%v')`, customer.Id, customer.Name))
}
