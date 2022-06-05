package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/marsel1323/YP-diploma-1/internal/config"
	"github.com/marsel1323/YP-diploma-1/internal/models"
	"github.com/marsel1323/YP-diploma-1/internal/repository"
	"github.com/marsel1323/YP-diploma-1/internal/utils"
	"log"
	"net/http"
)

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrLoginAndPasswordRequired = errors.New("login and Password are required")
var ErrInvalidJSON = errors.New("invalid json")

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
	// Parse data to json
	var json models.User

	if err := c.ShouldBindJSON(&json); err != nil {
		log.Println(ErrInvalidJSON)
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJSON})
		return
	}

	// Validate json
	if &json.Login != nil || &json.Password != nil {
		log.Println(ErrLoginAndPasswordRequired)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrLoginAndPasswordRequired})
		return
	}
	// Insert user in DB
	_, err := repo.DB.CreateUser(json)
	if errors.Is(err, ErrUserAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"error": "Login already used"})
	}

	sessionToken := utils.Hash(json.Login, repo.App.Config.Secret)
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
	// Parse data to json
	var json models.User
	if err := c.ShouldBindJSON(&json); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate json
	if &json.Login != nil || &json.Password != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login and Password are required"})
		return
	}

	// Check user in DB

	// Generate session token and set it to cookie
	sessionToken := utils.Hash(json.Login, repo.App.Config.Secret)
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
