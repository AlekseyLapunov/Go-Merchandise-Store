package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/handler"
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/usecase"
	_ "github.com/lib/pq"
)

const ENV_MERCH_STORE_DB_URL = "MERCH_STORE_DB_URL"
const ENV_MERCH_STORE_PORT   = "MERCH_STORE_PORT"

func init() {
    check := os.Getenv(ENV_MERCH_STORE_DB_URL)
    if check == "" {
        log.Fatalf("$\"%s\" not set: no DB connection URL", ENV_MERCH_STORE_DB_URL)
    }

    check = os.Getenv(ENV_MERCH_STORE_PORT)
    if check == "" {
        log.Fatalf("$\"%s\" not set: application port unknown", ENV_MERCH_STORE_PORT)
    }
}

func main() {
    dsn := os.Getenv(ENV_MERCH_STORE_DB_URL)

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Printf("Can't connect to the database: %v", err)
        return
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Printf("Probe ping to the database was not answered: %v", err)
        return
    }
    log.Println("Connection to the database is established")

    employeeStorage   := storage.NewEmployeeStorage(db)
    merchStorage      := storage.NewMerchStorage(db)
    managementStorage := storage.NewManagementStorage(db)

    employeeUsecase := usecase.NewEmployeeUsecase(&employeeStorage, &managementStorage)
    merchUsecase    := usecase.NewMerchUsecase(&merchStorage, &managementStorage)

    addr := fmt.Sprintf(":%s", os.Getenv(ENV_MERCH_STORE_PORT))
    router := handler.NewRouter(&employeeUsecase, &merchUsecase)
    server := &http.Server{
        Addr:    addr,
        Handler: router,
    }

    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Printf("Failed to start server: %v", err)
            return
        }
    }()
    log.Printf("Server started on %s\n", addr)

    <-done
    log.Println("Server is shutting down...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Printf("Server shutdown failed: %v", err)
        return
    }
    log.Println("Server stopped gracefully")
}
