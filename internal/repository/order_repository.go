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
	defer func(tx *sqlx.Tx) {
		err = tx.Rollback()
		if err != nil {
			log.Println("failed to rollback transaction:", err)
		}
	}(tx)

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
				item.Rid, item.Name, item.Sale, item.Size,
				item.TotalPrice, item.NmID, item.Brand, item.Status,
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
	return nil, nil
}
