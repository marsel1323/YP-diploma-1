package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/marsel1323/YP-diploma-1/internal/models"
	"github.com/marsel1323/YP-diploma-1/internal/utils"
	"log"
	"time"
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		//log.Fatalf("Unable to connect to database: %v", err)
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		//log.Fatalf("Unable to ping database: %v", err)
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return db, nil
}

type PostgresStorage struct {
	DB *sql.DB
}

func NewPostgresStorage(dsn string) (DBRepo, error) {
	log.Println("connecting to DB...")
	db, err := Connect(dsn)
	if err != nil {
		log.Println("error while connecting to DB...")
		return nil, err
	}
	log.Println("connected!!!")

	postgresStorage := &PostgresStorage{DB: db}

	return postgresStorage, nil
}

func (p *PostgresStorage) CreateUser(user models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	passwordHash, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	result, err := p.DB.ExecContext(
		ctx,
		`
  				INSERT INTO users(login, passwordhash) 
				VALUES ($1, $2);
			  `,
		user.Login,
		passwordHash,
	)
	if err != nil {
		return nil, err
	}
	user.ID, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &user, nil
}
