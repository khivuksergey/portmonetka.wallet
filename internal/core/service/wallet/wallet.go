package wallet

import (
	serviceerror "github.com/khivuksergey/portmonetka.wallet/error"
	"github.com/khivuksergey/portmonetka.wallet/internal/adapter/storage/entity"
	"github.com/khivuksergey/portmonetka.wallet/internal/core/port/repository"
	"github.com/khivuksergey/portmonetka.wallet/internal/core/port/service"
	"github.com/khivuksergey/portmonetka.wallet/internal/model"
	"strings"
)

type wallet struct {
	walletRepository repository.WalletRepository
}

func NewWalletService(repositoryManager *repository.Manager) service.WalletService {
	return &wallet{walletRepository: repositoryManager.Wallet}
}

func (w *wallet) GetWalletsByUserId(userId uint64) (*[]entity.Wallet, error) {
	return w.walletRepository.GetWalletsByUserId(userId)
}

func (w *wallet) CreateWallet(walletCreateDTO *model.WalletCreateDTO) (*entity.Wallet, error) {
	if w.walletRepository.ExistsWithName(walletCreateDTO.UserId, walletCreateDTO.Name) {
		return nil, serviceerror.WalletAlreadyExists
	}
	return w.walletRepository.CreateWallet(&entity.Wallet{
		UserId:        walletCreateDTO.UserId,
		Name:          walletCreateDTO.Name,
		Description:   walletCreateDTO.Description,
		Currency:      strings.ToUpper(walletCreateDTO.Currency),
		InitialAmount: walletCreateDTO.InitialAmount,
	})
}

func (w *wallet) UpdateWallet(walletUpdateDTO *model.WalletUpdateDTO) (*entity.Wallet, error) {
	walletToUpdate, err := w.walletRepository.GetWalletById(walletUpdateDTO.Id)
	if err != nil {
		return nil, serviceerror.WalletDoesntExist
	}
	err = w.validateUpdateWalletAttributes(walletToUpdate, walletUpdateDTO)
	if err != nil {
		return nil, err
	}
	return w.walletRepository.UpdateWallet(walletToUpdate)
}

func (w *wallet) DeleteWallet(walletDeleteDTO *model.WalletDeleteDTO) error {
	if !w.walletRepository.WalletBelongsToUser(walletDeleteDTO.Id, walletDeleteDTO.UserId) {
		return serviceerror.WalletDoesntBelongToUser
	}
	return w.walletRepository.DeleteWallet(walletDeleteDTO.Id)
}

// TODO move attributes validation to validator
func (w *wallet) validateUpdateWalletAttributes(wallet *entity.Wallet, walletUpdateDTO *model.WalletUpdateDTO) error {
	if walletUpdateDTO.Name == nil &&
		walletUpdateDTO.Description == nil &&
		walletUpdateDTO.Currency == nil &&
		walletUpdateDTO.InitialAmount == nil {
		return serviceerror.AtLeastOneFieldIsRequired
	}
	if walletUpdateDTO.Name != nil {
		if len(*walletUpdateDTO.Name) < 3 || len(*walletUpdateDTO.Name) > 128 {
			return serviceerror.WalletNameLengthError
		}
		if w.walletRepository.ExistsWithName(walletUpdateDTO.UserId, *walletUpdateDTO.Name) {
			return serviceerror.WalletAlreadyExists
		}
	}
	if walletUpdateDTO.Description != nil {
		if len(*walletUpdateDTO.Description) > 256 {
			return serviceerror.WalletDescriptionLengthError
		}
		wallet.Description = walletUpdateDTO.Description
	}
	if walletUpdateDTO.Currency != nil {
		if len(*walletUpdateDTO.Currency) != 3 {
			return serviceerror.WalletCurrencyError
		}
		wallet.Currency = *walletUpdateDTO.Currency
	}
	if walletUpdateDTO.InitialAmount != nil {
		wallet.InitialAmount = *walletUpdateDTO.InitialAmount
	}
	return nil
}
