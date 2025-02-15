package main

import (
    "context"
    "database/sql"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    "net/http"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/handler"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/usecase"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"
    _ "github.com/lib/pq"
)

func main() {
    const ENV_MERCH_STORE_DB_URL = "MERCH_STORE_DB_URL"

    dsn := os.Getenv(ENV_MERCH_STORE_DB_URL)
    if dsn == "" {
        log.Fatalf("\"%s\" environment variable is not set", ENV_MERCH_STORE_DB_URL)
    }

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatalf("Can't connect to the database: %v", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatalf("Probe ping to the database was not answered: %v", err)
    }
    log.Println("Connection to the database is established")

    employeeStorage   := storage.NewEmployeeStorage(db)
    merchStorage      := storage.NewMerchStorage(db)
    managementStorage := storage.NewManagementStorage(db)

    employeeUsecase := usecase.NewEmployeeUsecase(&employeeStorage, &managementStorage)
    merchUsecase    := usecase.NewMerchUsecase(merchStorage, managementStorage)

    router := handler.NewRouter(employeeUsecase, merchUsecase)
    server := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()
    log.Println("Server started on :8080")

    <-done
    log.Println("Server is shutting down...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server shutdown failed: %v", err)
    }
    log.Println("Server stopped gracefully")
}
