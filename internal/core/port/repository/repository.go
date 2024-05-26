package repository

import (
	"github.com/khivuksergey/portmonetka.wallet/internal/adapter/storage/entity"
)

type Manager struct {
	Wallet WalletRepository
}

//go:generate mockgen -source=repository.go -destination=../../../adapter/storage/gorm/repo/mock/mock_repository.go -package=mock
type WalletRepository interface {
	ExistsWithName(userId uint64, name string) bool
	WalletBelongsToUser(id, userId uint64) bool
	GetWalletById(id uint64) (*entity.Wallet, error)
	GetWalletsByUserId(userId uint64) ([]entity.Wallet, error)
	CreateWallet(wallet *entity.Wallet) (*entity.Wallet, error)
	UpdateWallet(wallet *entity.Wallet) (*entity.Wallet, error)
	DeleteWallet(id uint64) error
}
