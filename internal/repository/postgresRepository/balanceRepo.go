package postgresrepository

import (
	"context"
	"database/sql"
	"github.com/marsel1323/YP-diploma-1/internal/models"
	"log"
	"time"
)

func (p *PostgresStorage) SetBalance(userID int, value float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := p.DB.ExecContext(
		ctx,
		`INSERT INTO balances (current, user_id) 
			   VALUES ($1, $2)
			   ON CONFLICT (user_id)
			   DO UPDATE SET current = balances.current + $1;`,
		value,
		userID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStorage) GetBalance(userID int) (*models.Balance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := p.DB.QueryRowContext(
		ctx,
		`SELECT current, withdrawn
			   FROM balances 
			   WHERE user_id = $1
	   `,
		userID,
	)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var balance models.Balance
	err := row.Scan(
		&balance.Current,
		&balance.Withdrawn,
	)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (p *PostgresStorage) CreateWithdrawal(withdrawal *models.Withdrawal) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	_, err = p.DB.ExecContext(
		ctx,
		`UPDATE balances
				SET current = current - $1, withdrawn = withdrawn + $1
				WHERE user_id = $2;`,
		withdrawal.Sum,
		withdrawal.UserID,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = p.DB.ExecContext(
		ctx,
		`INSERT INTO withdrawals ("order", sum, processed_at, user_id)
				VALUES ($1, $2, now(), $3);`,
		withdrawal.Order,
		withdrawal.Sum,
		withdrawal.UserID,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (p *PostgresStorage) GetAllUserWithdrawals(userID int) ([]*models.Withdrawal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := p.DB.QueryContext(
		ctx,
		`SELECT "order", sum, processed_at
			   FROM withdrawals 
			   WHERE user_id = $1
			   ORDER BY processed_at;`,
		userID,
	)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(rows)
	if err != nil {
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	var withdrawals []*models.Withdrawal
	for rows.Next() {
		var withdrawal models.Withdrawal

		err := rows.Scan(
			&withdrawal.Order,
			&withdrawal.Sum,
			&withdrawal.ProcessedAt,
		)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, &withdrawal)
	}

	return withdrawals, nil
}
