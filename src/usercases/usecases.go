package usecases

import (
	"domain"
	"fmt"
)

type UserRepository interface {
	Store(user User)
	FindById(id int) User
}

type User struct {
	Id       int
	IsAdmin  bool
	Customer domain.Customer
}

type Item struct {
	Id    int
	Name  string
	value float64
}

type Logger interface {
	Log(message string) error
}

type OrderInteractor struct {
	UserRepository  UserRepository
	OrderRepository domain.OrderRepository
	ItemRepository  domain.ItemRepository
	Logger          Logger
}

func (interactor *OrderInteractor) Items(userId, orderId int) ([]Item, error) {
	var items []Item
	user := interactor.UserRepository.FindById(userId)
	order := interactor.OrderRepository.FindById(orderId)
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
		items[i] = Item{item.Id, item.Name, item.Value}
	}
	return items, nil
}

func (interactor *OrderInteractor) Add(userId, orderId, itemId int) error {
	return nil
}

type AdminOrderInteractor struct {
	OrderInteractor
}

func (interactor *AdminOrderInteractor) Add(userId, orderId, itemId int) error {
	return nil
}
