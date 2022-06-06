package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/marsel1323/YP-diploma-1/internal/config"
	"github.com/marsel1323/YP-diploma-1/internal/models"
	"github.com/marsel1323/YP-diploma-1/internal/repository"
	"io"
	"log"
	"net/http"
)

var ErrLoginAndPasswordRequired = errors.New("login and password are required")
var ErrInvalidRequest = errors.New("invalid request")
var ErrInvalidPassword = errors.New("invalid password")

type Repository struct {
	App *config.Application
	DB  repository.DBRepo
}

func (repo *Repository) GetAccrual(user *models.User, order *models.Order) {
	url := fmt.Sprintf(
		"%s/api/orders/%s",
		repo.App.Config.AccrualAddress,
		order.Number,
	)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("unable to get order", err)
		return
	}
	defer resp.Body.Close()

	log.Println("Accrual Response Status:", resp.StatusCode)
	order.Status = models.Processing
	err = repo.DB.UpdateOrder(order)
	if err != nil {
		log.Println("UpdateOrder:", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("io.ReadAll", err)
			return
		}

		respJSON := struct {
			Order   string            `json:"order"`
			Status  models.StatusType `json:"status"`
			Accrual float64           `json:"accrual"`
		}{}

		err = json.Unmarshal(bodyBytes, &respJSON)
		if err != nil {
			log.Println("json.Unmarshal", err)
			return
		}

		order.Status = respJSON.Status
		order.Accrual = respJSON.Accrual

		err = repo.DB.UpdateOrder(order)
		if err != nil {
			log.Println("UpdateOrder:", err)
			return
		}

		err = repo.DB.SetBalance(user.ID, order.Accrual)
		if err != nil {
			log.Println("SetBalance:", err)
			return
		}
	} else if resp.StatusCode == http.StatusTooManyRequests {
		log.Println("[Accrual]: Too many requests")
		return
	} else if resp.StatusCode == http.StatusInternalServerError {
		log.Println("[Accrual]: Internal server error")
		return
	}
}

func (repo *Repository) GetUser(c *gin.Context) (*models.User, error) {
	login, ok := c.Get("login")
	if !ok {
		return nil, errors.New("login not found")
	}

	loginStr, ok := login.(string)
	if !ok {
		return nil, errors.New("cannot type assert login")
	}

	user, err := repo.DB.GetUser(loginStr)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewRepo(appConfig *config.Application, db repository.DBRepo) *Repository {
	return &Repository{
		App: appConfig,
		DB:  db,
	}
}
