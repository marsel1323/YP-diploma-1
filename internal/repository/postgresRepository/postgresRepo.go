package postgresRepository

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/marsel1323/YP-diploma-1/internal/repository"
	"log"
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

func NewPostgresStorage(dsn string) (repository.DBRepo, error) {
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
