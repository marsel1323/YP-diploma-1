package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/marsel1323/YP-diploma-1/internal/models"
	"github.com/theplant/luhn"
	"log"
	"net/http"
	"strconv"
)

func (repo *Repository) GetBalance(c *gin.Context) {
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

	balance, err := repo.DB.GetBalance(user.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, balance)
}

func (repo *Repository) WithdrawBalance(c *gin.Context) {
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

	withdrawal := &models.Withdrawal{}

	if err := c.ShouldBindJSON(withdrawal); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidRequest})
		return
	}

	withdrawal.UserID = user.ID

	orderNumber, err := strconv.Atoi(withdrawal.Order)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	isValid := luhn.Valid(orderNumber)
	if !isValid {
		c.JSON(http.StatusUnprocessableEntity, nil)
		return
	}

	balance, err := repo.DB.GetBalance(withdrawal.UserID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	if balance.Current < withdrawal.Sum {
		log.Println("на счету недостаточно средств")
		c.JSON(http.StatusPaymentRequired, nil)
		return
	}

	err = repo.DB.CreateWithdrawal(withdrawal)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (repo *Repository) GetWithdrawalList(c *gin.Context) {
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

	withdrawals, err := repo.DB.GetAllUserWithdrawals(user.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	if len(withdrawals) == 0 {
		log.Println("204 — нет ни одного списания")
		c.JSON(http.StatusNoContent, nil)
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}
