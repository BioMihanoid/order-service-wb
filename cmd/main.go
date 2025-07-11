package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"order-service-wb/internal/api"
	"order-service-wb/internal/cache"
	"order-service-wb/internal/repository"
	"order-service-wb/internal/service"
	"order-service-wb/pkg/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
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

	repo := repository.NewOrderRepository(db)
	c := cache.NewCache(conf.Cache.Size)
	serv := service.NewOrderService(repo, c)

	if err = serv.LoadCache(context.Background(), conf.Cache.Size); err != nil {
		log.Fatalf("failed to load cache: %v", err)
	}

	handler := api.NewHandler(serv)

	server := &http.Server{
		Addr:    ":" + conf.Server.Port,
		Handler: handler.InitRouter(),
	}

	go func() {
		log.Println("Starting server on :" + conf.Server.Port)
		if err = server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("failed to gracefully shutdown server: %v", err)
	}
	log.Println("Server gracefully stopped")
}
