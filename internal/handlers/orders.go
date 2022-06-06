package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/marsel1323/YP-diploma-1/internal/models"
	"github.com/theplant/luhn"
	"io"
	"log"
	"net/http"
	"strconv"
)

func (repo *Repository) CreateOrder(c *gin.Context) {
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
		log.Println("GetUser", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	data, err := c.GetRawData()
	if err != nil {
		log.Println("c.GetRawData()", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidRequest})
		return
	}
	orderNumber, err := strconv.Atoi(string(data))
	if err != nil {
		log.Println("strconv.Atoi", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidRequest})
		return
	}

	isValid := luhn.Valid(orderNumber)
	if !isValid {
		c.JSON(http.StatusUnprocessableEntity, nil)
		return
	}

	order := &models.Order{
		Status: models.New,
		Number: strconv.Itoa(orderNumber),
	}

	existedOrder, err := repo.DB.GetOrder(order.Number)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		order, err = repo.DB.CreateOrder(user.ID, order)
		if err != nil {
			log.Println("CreateOrder", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
	} else if err != nil {
		log.Println("GetOrder", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	if existedOrder != nil {
		if existedOrder.UserID != user.ID {
			c.JSON(http.StatusConflict, nil)
			return
		}
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
	}()

	c.JSON(http.StatusAccepted, nil)
}

func (repo *Repository) GetAllOrders(c *gin.Context) {
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

	if len(orders) == 0 {
		c.SetCookie(
			"Content-Length",
			"0",
			0,
			"/",
			"localhost",
			false,
			true,
		)
		c.JSON(http.StatusNoContent, nil)
		return
	}
	c.JSON(http.StatusOK, orders)
}
