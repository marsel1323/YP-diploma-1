package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/marsel1323/YP-diploma-1/internal/models"
	"github.com/marsel1323/YP-diploma-1/internal/utils"
	"log"
	"net/http"
	"strings"
)

func (repo *Repository) RegisterUser(c *gin.Context) {
	var userJSON models.User

	if err := c.ShouldBindJSON(&userJSON); err != nil {
		log.Println(ErrInvalidRequest)
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidRequest})
		return
	}

	// Validate Login and Password
	if userJSON.Login == "" || userJSON.Password == "" {
		log.Println(ErrLoginAndPasswordRequired)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrLoginAndPasswordRequired})
		return
	}

	err := repo.DB.CreateUser(userJSON)

	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		log.Println(err)
		c.JSON(http.StatusConflict, gin.H{"error": "Login already used"})
		return
	} else if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	user, err := repo.DB.GetUser(userJSON.Login)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	//Create balance
	err = repo.DB.SetBalance(user.ID, 0)
	if err != nil {
		log.Println("Set balance")
		log.Println("Err", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	sessionToken := utils.Hash(user.Login, repo.App.Config.Secret)

	repo.App.Sessions[sessionToken] = user.Login

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
