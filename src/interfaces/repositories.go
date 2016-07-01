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

func (repo *DbOrderRepo) FindById(id int) domain.Order {
	row := repo.dbHandler.Query(fmt.Sprintf(`SELECT customer_id FROM orders 
                                            WHERE id = '%d'`,
		id))
	var customer_id int
	row.Next()
	row.Scan(&customer_id)
	dbCustomerRepo := NewDbCustomerRepo(repo.dbHandlers)
	customer := dbCustomerRepo.FindById(customer_id)
	order := domain.Order{Id: id, Customer: customer}
	var itemId int
	dbItemRepo := NewDbItemRepo(repo.dbHandlers)
	row = dbItemRepo.dbHandler.Query(fmt.Sprintf(`SELECT itme_id From items2order 
                                                    WHERE order_id = '%d'`,
		id))

	for row.Next() {
		row.Scan(&itemId)
		order.Add(dbItemRepo.FindById(itemId))
	}
	return order
}

func NewDbItemRepo(dbHandlers map[string]DbHandler) *DbItemRepo {
	repo := new(DbItemRepo)
	repo.dbHandlers = dbHandlers
	repo.dbHandler = dbHandlers["DbItemRepo"]
	return repo
}

func (repo *DbItemRepo) Store(item domain.Item) {
	available := "no"
	if item.Available {
		available = "yes"
	}
	repo.dbHandler.Execute(fmt.Sprintf(`INSERT INTO items (id, name, value, available) 
                                        VALUES ('%d', '%v', '%f', '%v')`,
		item.Id, item.Name, item.Value, available))
}

func (repo *DbItemRepo) FindById(id int) domain.Item {
	row := repo.dbHandler.Query(fmt.Sprintf(`SELECT name, value, available From items 
                                        WHERE id = 'id' LIMIT 1`,
		id))
	var name string
	var value float64
	var available string
	row.Next()
	row.Scan(&name, &value, &available)
	item := domain.Item{Id: id, Name: name, Value: value, Available: false}
	if available == "yes" {
		item.Available = true
	}
	return item
}
