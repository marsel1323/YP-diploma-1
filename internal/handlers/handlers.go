package handlers

import (
	"errors"
	"github.com/marsel1323/YP-diploma-1/internal/config"
	"github.com/marsel1323/YP-diploma-1/internal/repository"
)

var ErrLoginAndPasswordRequired = errors.New("login and password are required")
var ErrInvalidRequest = errors.New("invalid request")
var ErrInvalidPassword = errors.New("invalid password")

type Repository struct {
	App *config.Application
	DB  repository.DBRepo
}

func NewRepo(appConfig *config.Application, db repository.DBRepo) *Repository {
	return &Repository{
		App: appConfig,
		DB:  db,
	}
}
