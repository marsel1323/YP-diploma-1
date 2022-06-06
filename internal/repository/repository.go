package repository

import "github.com/marsel1323/YP-diploma-1/internal/models"

type DBRepo interface {
	CreateUser(user models.User) error
	GetUser(login string) (*models.User, error)
	CreateOrder(userID int, order *models.Order) (*models.Order, error)
	GetOrder(orderNumber string) (*models.Order, error)
	UpdateOrder(order *models.Order) error
	GetAllUserOrders(userID int) ([]*models.Order, error)
	GetBalance(userID int) (*models.Balance, error)
	SetBalance(userID int, value float64) error
	CreateWithdrawal(withdrawal *models.Withdrawal) error
	GetAllUserWithdrawals(userID int) ([]*models.Withdrawal, error)
}
