package repo

import (
	"github.com/khivuksergey/portmonetka.wallet/internal/adapter/storage/entity"
	"github.com/khivuksergey/portmonetka.wallet/internal/core/port/repository"
	"gorm.io/gorm"
)

type walletRepository struct {
	db        *gorm.DB
	tableName string
}

func NewWalletRepository(db *gorm.DB) repository.WalletRepository {
	return &walletRepository{db: db, tableName: entity.Wallet{}.TableName()}
}

func (w *walletRepository) ExistsWithName(userId uint64, name string) bool {
	var count int64
	w.db.Model(&entity.Wallet{}).Where("user_id = ? AND name = ?", userId, name).Count(&count)
	return count == 1
}

func (w *walletRepository) GetWalletById(id uint64) (*entity.Wallet, error) {
	wallet := &entity.Wallet{}
	result := w.db.First(wallet, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return wallet, nil
}

func (w *walletRepository) WalletBelongsToUser(id, userId uint64) bool {
	wallet, err := w.GetWalletById(id)
	if err != nil || wallet == nil {
		return false
	}
	return wallet.UserId == userId
}

func (w *walletRepository) GetWalletsByUserId(userId uint64) (*[]entity.Wallet, error) {
	wallets := &[]entity.Wallet{}
	result := w.db.
		Where("user_id = ?", userId).
		Order("updated_at desc").
		Find(wallets)
	if result.Error != nil {
		return nil, result.Error
	}
	return wallets, nil
}

func (w *walletRepository) CreateWallet(wallet *entity.Wallet) (*entity.Wallet, error) {
	if err := w.db.Create(wallet).Error; err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w *walletRepository) UpdateWallet(wallet *entity.Wallet) (*entity.Wallet, error) {
	err := w.db.Save(wallet).Error
	return wallet, err
}

func (w *walletRepository) DeleteWallet(id uint64) error {
	return w.db.Delete(&entity.Wallet{}, id).Error
}
