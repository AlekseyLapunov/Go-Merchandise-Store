package usecase

import "github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"

type MerchUsecase struct {
	storage storage.MerchStorage
}

func NewMerchUsecase(s storage.MerchStorage) MerchUsecase {
	return MerchUsecase{storage: s}
}
