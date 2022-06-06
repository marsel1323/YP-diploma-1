package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/marsel1323/YP-diploma-1/internal/config"
	"github.com/marsel1323/YP-diploma-1/internal/models"
	"github.com/marsel1323/YP-diploma-1/internal/repository"
	"github.com/marsel1323/YP-diploma-1/internal/utils"
	"github.com/theplant/luhn"
	"io"
	"log"
	"net/http"
	"strconv"
)

var ErrUserAlreadyExists = errors.New("user already exists")
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

func (repo *Repository) RegisterUser(c *gin.Context) {
	var userJSON models.User

	if err := c.ShouldBindJSON(&userJSON); err != nil {
		log.Println(ErrInvalidRequest)
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidRequest})
		return
	}

	if userJSON.Login == "" || userJSON.Password == "" {
		log.Println(ErrLoginAndPasswordRequired)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrLoginAndPasswordRequired})
		return
	}

	user, err := repo.DB.GetUser(userJSON.Login)

	if !errors.Is(err, sql.ErrNoRows) {
		log.Println("Err", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	if user != nil {
		if user.Login == userJSON.Login {
			c.JSON(http.StatusConflict, gin.H{"error": "login already in use"})
			return
		}
	}

	_, err = repo.DB.CreateUser(userJSON)
	if errors.Is(err, ErrUserAlreadyExists) {
		log.Println(err)
		c.JSON(http.StatusConflict, gin.H{"error": "Login already used"})
		return
	}

	sessionToken := utils.Hash(userJSON.Login, repo.App.Config.Secret)

	repo.App.Sessions[sessionToken] = userJSON.Login

	c.SetCookie(
		"Authorization",
		sessionToken,
		24*60*60,
		"/",
		"localhost",
		false,
		true,
	)
	c.JSON(http.StatusOK, nil)
}

func (repo *Repository) LoginUser(c *gin.Context) {
	var userJSON models.User

	if err := c.ShouldBindJSON(&userJSON); err != nil {
		log.Println(ErrInvalidRequest)
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidRequest})
		return
	}

	if userJSON.Login == "" || userJSON.Password == "" {
		log.Println(ErrLoginAndPasswordRequired)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrLoginAndPasswordRequired})
		return
	}

	user, err := repo.DB.GetUser(userJSON.Login)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	passwordIsValid := utils.ComparePassword(userJSON.Password, user.Password)
	if !passwordIsValid {
		log.Println(ErrInvalidPassword)
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidPassword})
		return
	}

	sessionToken := utils.Hash(userJSON.Login, repo.App.Config.Secret)

	repo.App.Sessions[sessionToken] = userJSON.Login

	c.SetCookie(
		"Authorization",
		sessionToken,
		24*60*60,
		"/",
		"localhost",
		false,
		true,
	)
	c.JSON(http.StatusOK, nil)
}

func (repo *Repository) CreateOrder(c *gin.Context) {
	order := &models.Order{
		Status: models.New,
	}

	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidRequest})
		return
	}
	number, err := strconv.Atoi(string(data))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidRequest})
		return
	}

	// 422 — неверный формат номера заказа
	isValid := luhn.Valid(number)
	if !isValid {
		c.JSON(http.StatusUnprocessableEntity, nil)
		return
	}

	login, ok := c.Get("login")
	if !ok {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}
	loginStr, ok := login.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	user, err := repo.DB.GetUser(loginStr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	//log.Println(user)

	order.Number = strconv.Itoa(number)

	existedOrder, err := repo.DB.GetOrder(order.Number)
	if errors.Is(err, sql.ErrNoRows) {
		order, err = repo.DB.CreateOrder(user.ID, order)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
	} else if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	if existedOrder != nil {
		// 409 — номер заказа уже был загружен другим пользователем
		if existedOrder.UserID != user.ID {
			log.Println("номер заказа уже был загружен другим пользователем")
			c.JSON(http.StatusConflict, nil)
			return
		}
		// 200 — номер заказа уже был загружен этим пользователем;
		log.Println("номер заказа уже был загружен этим пользователем")
		c.JSON(http.StatusOK, nil)
		return
	}

	go func() {
		url := fmt.Sprintf(
			"%s/api/orders/%s",
			repo.App.Config.AccrualAddress,
			order.Number,
		)

		resp, err := http.Get(url)
		if err != nil {
			log.Println("unable to get order")
		}
		defer resp.Body.Close()

		//log.Println("Response status:", resp.StatusCode)
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				return
			}

			respJSON := struct {
				Order   string            `json:"order"`
				Status  models.StatusType `json:"status"`
				Accrual int               `json:"accrual"`
			}{}

			err = json.Unmarshal(bodyBytes, &respJSON)
			if err != nil {
				log.Println(err)
				return
			}

			order.Status = respJSON.Status
			order.Accrual = respJSON.Accrual

			err = repo.DB.UpdateOrder(order)
			if err != nil {
				log.Println(err)
				return
			}
		} else if resp.StatusCode == http.StatusTooManyRequests {

		} else if resp.StatusCode == http.StatusInternalServerError {

		}

		//err = resp.Body.Close()
		//if err != nil {
		//	log.Println(err)
		//}
	}()
	// 202 — новый номер заказа принят в обработку
	c.JSON(http.StatusAccepted, nil)
}

func (repo *Repository) GetAllOrders(c *gin.Context) {
	// 200 — успешная обработка запроса.
	login, ok := c.Get("login")
	if !ok {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}
	loginStr, ok := login.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	user, err := repo.DB.GetUser(loginStr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	orders, err := repo.DB.GetAllUserOrders(user.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	// 204 — нет данных для ответа.
	if len(orders) == 0 {
		log.Println("204 — нет данных для ответа.")
		c.JSON(http.StatusNoContent, nil)
		return
	}
	c.JSON(http.StatusOK, orders)
}
