package postgresRepository

import (
	"context"
	"database/sql"
	"github.com/marsel1323/YP-diploma-1/internal/models"
	"log"
	"time"
)

func (p *PostgresStorage) CreateOrder(userID int, order *models.Order) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println(userID, order)
	_, err := p.DB.ExecContext(
		ctx,
		`INSERT INTO orders(number, status, accrual, uploaded_at, user_id)
			   VALUES ($1, $2, $3, $4, $5);`,
		order.Number,
		order.Status,
		order.Accrual,
		time.Now(),
		userID,
	)
	if err != nil {
		return nil, err
	}

	//orderID, err := result.LastInsertId()
	//if err != nil {
	//	return nil, err
	//}
	//order.ID = int(orderID)

	return order, nil
}

func (p *PostgresStorage) GetOrder(orderNumber string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := p.DB.QueryRowContext(
		ctx,
		`SELECT id, number, status, accrual, uploaded_at, user_id
			   FROM orders 
			   WHERE number = $1
	   `,
		orderNumber,
	)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var o models.Order
	err := row.Scan(
		&o.ID,
		&o.Number,
		&o.Status,
		&o.Accrual,
		&o.UploadedAt,
		&o.UserID,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (p *PostgresStorage) UpdateOrder(order *models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println(order)

	_, err := p.DB.ExecContext(
		ctx,
		`UPDATE orders SET status = $1, accrual = $2 WHERE number = $3`,
		order.Status,
		order.Accrual,
		order.Number,
	)
	if err != nil {
		return err
	}

	//orderID, err := result.LastInsertId()
	//if err != nil {
	//	return nil, err
	//}
	//order.ID = int(orderID)

	return nil
}

func (p *PostgresStorage) GetAllUserOrders(userID int) ([]*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := p.DB.QueryContext(
		ctx,
		`SELECT number, status, accrual, uploaded_at
			   FROM orders 
			   WHERE user_id = $1
			   ORDER BY uploaded_at;`,
		userID,
	)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	if err != nil {
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	var orders []*models.Order
	for rows.Next() {
		var order models.Order

		err := rows.Scan(
			//&order.ID,
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
			//&order.UserID,
		)
		if err != nil {
			return nil, err
		}

		//
		//uploadedAt, err := time.Parse(time.RFC3339, order.UploadedAt)
		//order.UploadedAt = uploadedAt.String()
		//order.UploadedAt = order.UploadedAt.Format(time.RFC3339)
		orders = append(orders, &order)
	}

	return orders, nil
}
