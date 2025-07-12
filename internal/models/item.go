package models

type Item struct {
	ChrtID      int    `json:"chrt_id" db:"chrt_id" validate:"required,gte=0"`
	TrackNumber string `json:"track_number" db:"track_number" validate:"required"`
	Price       int    `json:"price" db:"price" validate:"required,gte=0"`
	Rid         string `json:"rid" db:"rid" validate:"required"`
	Name        string `json:"name" db:"name" validate:"required"`
	Sale        int    `json:"sale" db:"sale" validate:"gte=0"`
	Size        string `json:"size" db:"size" validate:"required"`
	TotalPrice  int    `json:"total_price" db:"total_price" validate:"required,gte=0"`
	NmID        int    `json:"nm_id" db:"nm_id" validate:"required,gte=0"`
	Brand       string `json:"brand" db:"brand" validate:"required"`
	Status      int    `json:"status" db:"status" validate:"gte=0"`
}
