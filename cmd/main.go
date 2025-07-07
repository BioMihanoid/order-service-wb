package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"order-service-wb/internal/repository"
	"order-service-wb/pkg/config"
)

func main() {
	conf := config.NewConfig()
	db, err := sqlx.Connect("postgres", conf.DbConfig.GetDSN())
	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			log.Fatalf("failed to close database connection: %v", err)
		}
	}(db)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	_ = repository.NewOrderRepository(db)
}
