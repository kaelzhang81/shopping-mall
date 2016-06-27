package usecases

import (
	"domain"
	"fmt"
)

type UserRespository interface {
	Store(user User)
	GetById(id int) User
}

type User struct {
	Id       int
	IsAdmin  bool
	customer domain.Customer
}

type Logger interface {
	Log(message string) error
}

type OrderInteractor struct {
	UserRespository  UserRespository
	OrderRespository domain.OrderRespository
	ItemRespository  domain.ItemRespository
	Logger           Logger
}

func (interactor *OrderInteractor) Items(userId, orderId int) ([]Item, error) {
	var items []Item
	user := UserRespository.FindById(userId)
	order := OrderRespository.FindById(orderId)
	if user.Customer.Id != order.Customer.Id {
		message := "User #%i (customer #%) "
		message += "is not allowed to see items "
		message += "in order #%i (of customer #%i)"
		err := fmt.Errorf(message,
			userId,
			user.Customer.Id,
			orderId,
			order.Customer.Id)
		items = make([]Item, 0)
		return items, err
	}

	items = make([]Item, len(order.Items))
	for i, item := range order.Items {
		items[i] = Item(item.Id, item.Name, item.Value)
	}
}

func (interactor *OrderInteractor) Add(userId, orderId, itemId int) error {

}

type AdminOrderInteractor struct {
	OrderInteractor
}

func (interactor *AdminOrderInteractor) Add(userId, orderId, itemId int) error {

}
