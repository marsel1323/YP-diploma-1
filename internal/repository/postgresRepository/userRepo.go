package postgresrepository

import (
	"context"
	"github.com/marsel1323/YP-diploma-1/internal/models"
	"github.com/marsel1323/YP-diploma-1/internal/utils"
	"time"
)

func (p *PostgresStorage) CreateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	passwordHash, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = p.DB.ExecContext(
		ctx,
		`
  				INSERT INTO users(login, passwordhash) 
				VALUES ($1, $2);
			  `,
		user.Login,
		passwordHash,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStorage) GetUser(login string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := p.DB.QueryRowContext(
		ctx,
		`SELECT id, login, passwordhash
				FROM users
				WHERE login = $1`,
		login,
	)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.Login,
		&u.Password,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
