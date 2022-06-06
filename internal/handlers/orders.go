package handlers

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/marsel1323/YP-diploma-1/internal/models"
	"github.com/theplant/luhn"
	"log"
	"net/http"
	"strconv"
)

func (repo *Repository) CreateOrder(c *gin.Context) {
	user, err := repo.GetUser(c)
	if err != nil {
		log.Println(err)
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

	go repo.GetAccrual(user, order)

	c.JSON(http.StatusAccepted, nil)
}

func (repo *Repository) GetAllOrders(c *gin.Context) {
	user, err := repo.GetUser(c)
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
