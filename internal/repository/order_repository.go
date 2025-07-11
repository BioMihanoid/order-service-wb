package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"order-service-wb/internal/models"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
	GetAllOrders(ctx context.Context, limit int) ([]*models.Order, error)
}

type orderRepo struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepo{
		db: db,
	}
}

func (r *orderRepo) CreateOrder(ctx context.Context, order *models.Order) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println("failed to begin transaction:", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = ctx.Err(); err != nil {
		log.Println("context error before execution:", err)
		return fmt.Errorf("context cancelled before execution: %w", err)
	}

	q := `INSERT INTO orders(
            order_uid, track_number, entry, locale, 
        	internal_signature, customer_id, delivery_service,
        	shardkey, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`

	_, err = tx.ExecContext(ctx, q,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSig, order.CustomerID, order.DeliverySrv,
		order.ShardKey, order.SmID, order.DateCreated, order.OofShard,
	)

	if err != nil {
		log.Println("failed to execute insert order query:", err)
		return fmt.Errorf("failed to insert order: %w", err)
	}

	if err = ctx.Err(); err != nil {
		log.Println("context error before execution:", err)
		return fmt.Errorf("context cancelled before execution: %w", err)
	}

	q = `INSERT INTO items(order_uid, chrt_id, 
                  track_number, price, rid, name, sale, 
                  size, total_price, nm_id, brand, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $9, $8, $10, $11, $12)
		`

	for _, item := range order.Items {

		select {
		case <-ctx.Done():
			log.Println("context cancelled before executing items insert:", ctx.Err())
			return fmt.Errorf("context cancelled before executing items insert: %w", ctx.Err())
		default:
			_, err = tx.ExecContext(ctx, q,
				order.OrderUID, item.ChrtID, item.TrackNumber, item.Price,
				item.Rid, item.Name, item.Sale, item.TotalPrice, item.Size,
				item.NmID, item.Brand, item.Status,
			)
			if err != nil {
				log.Println("failed to execute insert items query:", err)
				return fmt.Errorf("failed to insert item: %w", err)
			}
		}
	}

	q = `INSERT INTO payment(order_uid, transaction, request_id, currency, 
                    provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`
	_, err = tx.ExecContext(ctx, q,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee,
	)

	if err != nil {
		log.Println("failed to execute insert payment query:", err)
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	if err = ctx.Err(); err != nil {
		log.Println("context error before execution:", err)
		return fmt.Errorf("context cancelled before execution: %w", err)
	}

	q = `INSERT INTO delivery(order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`

	_, err = tx.ExecContext(ctx, q,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone,
		order.Delivery.Zip, order.Delivery.City, order.Delivery.Addr,
		order.Delivery.Region, order.Delivery.Email,
	)

	if err != nil {
		log.Println("failed to execute insert delivery query:", err)
		return fmt.Errorf("failed to insert delivery: %w", err)
	}

	if err = ctx.Err(); err != nil {
		log.Println("context error before execution:", err)
		return fmt.Errorf("context cancelled before execution: %w", err)
	}

	if err = tx.Commit(); err != nil {
		log.Println("failed to commit transaction:", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *orderRepo) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println("failed to begin transaction:", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *sqlx.Tx) {
		err = tx.Rollback()
		if err != nil {
			log.Println("failed to rollback transaction:", err)
		}
	}(tx)

	if err = ctx.Err(); err != nil {
		log.Println("context error before execution:", err)
		return nil, fmt.Errorf("context cancelled before execution: %w", err)
	}

	var order models.Order
	q := `SELECT 
			order_uid, track_number, entry, locale,
			internal_signature, customer_id, delivery_service,
			shardkey, sm_id, date_created, oof_shard
		FROM orders WHERE order_uid = $1
		`
	err = tx.GetContext(ctx, &order, q, orderID)
	if err != nil {
		log.Println("failed to get order by ID:", err)
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
	}

	if err = ctx.Err(); err != nil {
		log.Println("context error after getting order:", err)
		return nil, fmt.Errorf("context cancelled after getting order: %w", err)
	}

	q = `SELECT
			chrt_id, track_number, price, rid, name, sale,
			size, total_price, nm_id, brand, status
		FROM items WHERE order_uid = $1
		`
	err = tx.SelectContext(ctx, &order.Items, q, orderID)
	if err != nil {
		log.Println("failed to get items for order:", err)
		return nil, fmt.Errorf("failed to get items for order: %w", err)
	}

	if err = ctx.Err(); err != nil {
		log.Println("context error after getting items:", err)
		return nil, fmt.Errorf("context cancelled after getting items: %w", err)
	}

	q = `SELECT
			transaction, request_id, currency, provider,
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment WHERE order_uid = $1
		`
	err = tx.GetContext(ctx, &order.Payment, q, orderID)
	if err != nil {
		log.Println("failed to get payment for order:", err)
		return nil, fmt.Errorf("failed to get payment for order: %w", err)
	}

	if err = ctx.Err(); err != nil {
		log.Println("context error after getting payment:", err)
		return nil, fmt.Errorf("context cancelled after getting payment: %w", err)
	}

	q = `SELECT
			name, phone, zip, city, address, region, email
		FROM delivery WHERE order_uid = $1
		`
	err = tx.GetContext(ctx, &order.Delivery, q, orderID)
	if err != nil {
		log.Println("failed to get delivery for order:", err)
		return nil, fmt.Errorf("failed to get delivery for order: %w", err)
	}

	if err = ctx.Err(); err != nil {
		log.Println("context error after getting delivery:", err)
		return nil, fmt.Errorf("context cancelled after getting delivery: %w", err)
	}

	if err = tx.Commit(); err != nil {
		log.Println("failed to commit transaction:", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &order, nil
}

func (r *orderRepo) GetAllOrders(ctx context.Context, limit int) ([]*models.Order, error) {
	q := `SELECT order_uid FROM orders ORDER BY date_created DESC LIMIT $1`

	rows, err := r.db.QueryxContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			continue
		}
		order, err := r.GetOrderByID(ctx, id)
		if err == nil {
			orders = append(orders, order)
		}
	}

	return orders, nil
}
