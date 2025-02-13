package main

import (
    "database/sql"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/handler"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/usecase"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"
)

func main() {
    db, err := sql.Open("", "")

    if err != nil {
        // TODO graceful shutdown
    }
    defer db.Close()

    employeeStorage := storage.NewEmployeeStorage(db)
    merchStorage    := storage.NewMerchStorage(db)

    employeeUsecase := usecase.NewEmployeeUsecase(employeeStorage)
    merchUsecase    := usecase.NewMerchUsecase(merchStorage)

    router := handler.NewRouter(employeeUsecase, merchUsecase)

    if err := router.Run(":8080"); err != nil {
        // TODO graceful shutdown
    }
}
