package repository

import "github.com/marsel1323/YP-diploma-1/internal/models"

type DBRepo interface {
	CreateUser(models.User) (*models.User, error)
}
