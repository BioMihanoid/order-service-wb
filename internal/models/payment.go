package models

type Payment struct {
	Transaction  string `json:"transaction" db:"transaction" validate:"required"`
	RequestID    string `json:"request_id" db:"request_id" validate:"required"`
	Currency     string `json:"currency" db:"currency" validate:"required,len=3"`
	Provider     string `json:"provider" db:"provider" validate:"required"`
	Amount       int    `json:"amount" db:"amount" validate:"required,gte=0"`
	PaymentDT    int64  `json:"payment_dt" db:"payment_dt" validate:"required,gte=0"`
	Bank         string `json:"bank" db:"bank" validate:"required"`
	DeliveryCost int    `json:"delivery_cost" db:"delivery_cost" validate:"gte=0"`
	GoodsTotal   int    `json:"goods_total" db:"goods_total" validate:"gte=0"`
	CustomFee    int    `json:"custom_fee" db:"custom_fee" validate:"gte=0"`
}
