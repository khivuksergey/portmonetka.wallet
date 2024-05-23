package model

import "github.com/shopspring/decimal"

type WalletCreateDTO struct {
	UserId        uint64          `json:"userId"`
	Name          string          `json:"name" validate:"required"`
	Description   *string         `json:"description"`
	Currency      string          `json:"currency" validate:"required,len=3"`
	InitialAmount decimal.Decimal `json:"initialAmount"`
}

type WalletUpdateDTO struct {
	Id            uint64           `json:"id"`
	UserId        uint64           `json:"userId"`
	Name          *string          `json:"name"`
	Description   *string          `json:"description"`
	Currency      *string          `json:"currency"`
	InitialAmount *decimal.Decimal `json:"initialAmount"`
}

type WalletDeleteDTO struct {
	Id     uint64 `json:"id"`
	UserId uint64 `json:"userId"`
}
