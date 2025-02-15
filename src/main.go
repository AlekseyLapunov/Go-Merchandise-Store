package main

import (
    "database/sql"
    "log"
    "os"
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

    employeeUsecase := usecase.NewEmployeeUsecase(employeeStorage, managementStorage)
    merchUsecase    := usecase.NewMerchUsecase(merchStorage, managementStorage)

    router := handler.NewRouter(employeeUsecase, merchUsecase)

    if err := router.Run(":8080"); err != nil {
        // TODO graceful shutdown
    }
}
