package models

import "time"

type Order struct {
	OrderUID    string    `json:"order_uid" db:"order_uid" validate:"required"`
	TrackNumber string    `json:"track_number" db:"track_number" validate:"required"`
	Entry       string    `json:"entry" db:"entry" validate:"required"`
	Delivery    Delivery  `json:"delivery" db:"-" validate:"required"`
	Payment     Payment   `json:"payment" db:"-" validate:"required"`
	Items       []Item    `json:"items" db:"-" validate:"required,min=1,dive"`
	Locale      string    `json:"locale" db:"locale" validate:"required"`
	InternalSig string    `json:"internal_signature" db:"internal_signature"`
	CustomerID  string    `json:"customer_id" db:"customer_id" validate:"required"`
	DeliverySrv string    `json:"delivery_service" db:"delivery_service"`
	ShardKey    string    `json:"shardkey" db:"shardkey"`
	SmID        int       `json:"sm_id" db:"sm_id" validate:"gte=0"`
	DateCreated time.Time `json:"date_created" db:"date_created" validate:"required"`
	OofShard    string    `json:"oof_shard" db:"oof_shard"`
}
