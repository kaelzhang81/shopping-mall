package interfaces

import (
	"domain"
	"fmt"
	"usecases"
)

type DbHandler interface {
	Execute(statment string)
	Query(statement string) Row
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

func (repo *DbUserRepo) FindById(id int) usecases.User {
	row := repo.dbHandler.Query(fmt.Sprintf(`SELECT is_admin, customer_id FROM users WHERE id = '%d' LIMIT 1`, id))
	var isAdmin string
	var customerId int
	row.Next()
	row.Scan(&isAdmin, &customerId)
	customerRepo := NewDbCustomerRepo(repo.dbHandlers)
	user := usecases.User{Id: id, Customer: customerRepo.FindById(customerId)}
	user.IsAdmin = false
	if isAdmin == "yes" {
		user.IsAdmin = true
	}
	return user
}

func NewDbCustomerRepo(dbHandlers map[string]DbHandler) *DbCustomerRepo {
	repo := new(DbCustomerRepo)
	repo.dbHandlers = dbHandlers
	repo.dbHandler = dbHandlers["DbCustomerRepo"]
	return repo
}

func (repo *DbCustomerRepo) Store(customer domain.Customer) {
	repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO customers (id, name) 
                                        VALUES ('%d', '%v')`, customer.Id, customer.Name))
}

func (repo *DbCustomerRepo) FindById(id int) domain.Customer {
	row := repo.dbHandler.Query(fmt.Sprintf(`SELECT name FROM customers WHERE id = '%d' LIMIT 1`, id))
	var name string
	row.Next()
	row.Scan(&name)
	customer := domain.Customer{Id: id, Name: name}
	return customer
}

func NewDbOrderRepo(dbHandlers map[string]DbHandler) *DbOrderRepo {
	repo := new(DbOrderRepo)
	repo.dbHandlers = dbHandlers
	repo.dbHandler = dbHandlers["DbOrderRepo"]
	return repo
}

func (repo *DbOrderRepo) Store(order domain.Order) {
	repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO orders (id, customer_id) 
                                        VALUES ('%d', '%v')`,
		order.Id, order.Customer.Id))

	for _, item := range order.Items {
		repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO items2order (item_id, order_id)
                                            VALUES ('%d', '%d')`,
			item.Id, order.Id))
	}
}
