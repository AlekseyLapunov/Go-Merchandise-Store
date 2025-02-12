package main

import (
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/handler"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/usecase"
)

func main() {

    employeeUsecase := usecase.NewEmployeeUsecase(/*storage*/)
    merchUsecase    := usecase.NewMerchUsecase(/*storage*/)

    router := handler.NewRouter(employeeUsecase, merchUsecase)

    if err := router.Run(":8080"); err != nil {
        // TODO graceful shutdown
    }
}
