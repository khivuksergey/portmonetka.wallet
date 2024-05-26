package service

import (
	"github.com/khivuksergey/portmonetka.wallet/internal/adapter/storage/entity"
	"github.com/khivuksergey/portmonetka.wallet/internal/model"
)

type Manager struct {
	Wallet WalletService
}

type WalletService interface {
	GetWalletsByUserId(userId uint64) ([]entity.Wallet, error)
	CreateWallet(walletCreateDTO model.WalletCreateDTO) (*entity.Wallet, error)
	UpdateWallet(walletUpdateDTO model.WalletUpdateDTO) (*entity.Wallet, error)
	DeleteWallet(walletDeleteDTO model.WalletDeleteDTO) error
}
