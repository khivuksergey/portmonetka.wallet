package wallet

import (
	serviceerror "github.com/khivuksergey/portmonetka.wallet/error"
	"github.com/khivuksergey/portmonetka.wallet/internal/adapter/storage/entity"
	"github.com/khivuksergey/portmonetka.wallet/internal/adapter/storage/gorm/repo/mock"
	"github.com/khivuksergey/portmonetka.wallet/internal/core/port/repository"
	"github.com/khivuksergey/portmonetka.wallet/internal/core/service/wallet"
	"github.com/khivuksergey/portmonetka.wallet/internal/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"strings"
	"testing"
)

func TestGetWalletsByUserId_Success(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWalletRepository := mock.NewMockWalletRepository(ctl)
	mockManager := &repository.Manager{
		Wallet: mockWalletRepository,
	}

	walletService := wallet.NewWalletService(mockManager)

	userId := uint64(1)
	expectedWallets := &[]entity.Wallet{
		{
			Id:            1,
			UserId:        userId,
			Name:          "Test Wallet 1",
			Description:   stringPtr("Description 1"),
			Currency:      "USD",
			InitialAmount: decimal.NewFromFloat(100.00),
		},
		{
			Id:            2,
			UserId:        userId,
			Name:          "Test Wallet 2",
			Description:   stringPtr("Description 2"),
			Currency:      "EUR",
			InitialAmount: decimal.NewFromFloat(200.00),
		},
	}

	mockWalletRepository.
		EXPECT().
		GetWalletsByUserId(userId).
		Times(1).
		Return(expectedWallets, nil)

	actualWallets, err := walletService.GetWalletsByUserId(userId)

	assert.NoError(t, err)
	assert.NotNil(t, actualWallets)
	assert.Equal(t, expectedWallets, actualWallets)
}

func TestCreateWallet_Success(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWalletRepository := mock.NewMockWalletRepository(ctl)
	mockManager := &repository.Manager{
		Wallet: mockWalletRepository,
	}

	walletService := wallet.NewWalletService(mockManager)

	walletCreateDTO := &model.WalletCreateDTO{
		UserId:        1,
		Name:          "Test wallet",
		Description:   stringPtr("A test wallet"),
		Currency:      "usd",
		InitialAmount: decimal.NewFromFloat(123.45),
	}

	expectedWallet := &entity.Wallet{
		UserId:        walletCreateDTO.UserId,
		Name:          walletCreateDTO.Name,
		Description:   walletCreateDTO.Description,
		Currency:      strings.ToUpper(walletCreateDTO.Currency),
		InitialAmount: walletCreateDTO.InitialAmount,
	}

	mockWalletRepository.
		EXPECT().
		ExistsWithName(walletCreateDTO.UserId, walletCreateDTO.Name).
		Times(1).
		Return(false)

	mockWalletRepository.
		EXPECT().
		CreateWallet(expectedWallet).
		Times(1).
		Return(expectedWallet, nil)

	createdWallet, err := walletService.CreateWallet(walletCreateDTO)

	assert.NoError(t, err)
	assert.NotNil(t, createdWallet)
	assert.Equal(t, createdWallet, expectedWallet)
}

func TestCreateWallet_DuplicateName_Error(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWalletRepository := mock.NewMockWalletRepository(ctl)
	mockManager := &repository.Manager{
		Wallet: mockWalletRepository,
	}

	walletService := wallet.NewWalletService(mockManager)

	walletCreateDTO := &model.WalletCreateDTO{
		UserId:        1,
		Name:          "Duplicate wallet",
		Description:   stringPtr("A test wallet with duplicate name"),
		Currency:      "USD",
		InitialAmount: decimal.NewFromFloat(123.45),
	}

	mockWalletRepository.
		EXPECT().
		ExistsWithName(walletCreateDTO.UserId, walletCreateDTO.Name).
		Times(1).
		Return(true)

	createdWallet, err := walletService.CreateWallet(walletCreateDTO)

	assert.Error(t, err)
	assert.Nil(t, createdWallet)
	assert.Equal(t, serviceerror.WalletAlreadyExists, err)
}

func TestUpdateWallet_Success(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWalletRepository := mock.NewMockWalletRepository(ctl)
	mockManager := &repository.Manager{
		Wallet: mockWalletRepository,
	}

	walletService := wallet.NewWalletService(mockManager)

	walletUpdateDTO := &model.WalletUpdateDTO{
		Id:            1,
		UserId:        1,
		Name:          stringPtr("Updated wallet name"),
		Description:   stringPtr("Updated description"),
		Currency:      stringPtr("eur"),
		InitialAmount: decimalPtr(decimal.NewFromFloat(150.00)),
	}

	existingWallet := &entity.Wallet{
		Id:            1,
		UserId:        1,
		Name:          "Old wallet name",
		Description:   stringPtr("Old description"),
		Currency:      "USD",
		InitialAmount: decimal.NewFromFloat(100.00),
	}

	updatedWallet := &entity.Wallet{
		Id:            1,
		UserId:        1,
		Name:          "Updated wallet name",
		Description:   stringPtr("Updated description"),
		Currency:      "EUR",
		InitialAmount: decimal.NewFromFloat(150.00),
	}

	mockWalletRepository.
		EXPECT().
		GetWalletById(walletUpdateDTO.Id).
		Times(1).
		Return(existingWallet, nil)

	mockWalletRepository.
		EXPECT().
		ExistsWithName(existingWallet.UserId, *walletUpdateDTO.Name).
		Times(1).
		Return(false)

	mockWalletRepository.
		EXPECT().
		UpdateWallet(existingWallet).
		Times(1).
		DoAndReturn(func(wallet *entity.Wallet) (*entity.Wallet, error) {
			wallet.Name = *walletUpdateDTO.Name
			wallet.Description = walletUpdateDTO.Description
			wallet.Currency = strings.ToUpper(*walletUpdateDTO.Currency)
			wallet.InitialAmount = *walletUpdateDTO.InitialAmount
			return wallet, nil
		})

	updatedWalletFromService, err := walletService.UpdateWallet(walletUpdateDTO)

	assert.NoError(t, err)
	assert.NotNil(t, updatedWalletFromService)
	assert.Equal(t, updatedWallet, updatedWalletFromService)
}

func TestUpdateWallet_WalletNotFound_Error(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWalletRepository := mock.NewMockWalletRepository(ctl)
	mockManager := &repository.Manager{
		Wallet: mockWalletRepository,
	}

	walletService := wallet.NewWalletService(mockManager)

	walletUpdateDTO := &model.WalletUpdateDTO{
		Id:            1,
		Name:          stringPtr("Non-existent wallet"),
		Description:   stringPtr("This wallet does not exist"),
		Currency:      stringPtr("USD"),
		InitialAmount: decimalPtr(decimal.NewFromFloat(200.00)),
	}

	mockWalletRepository.
		EXPECT().
		GetWalletById(walletUpdateDTO.Id).
		Times(1).
		Return(nil, serviceerror.WalletDoesntExist)

	updatedWallet, err := walletService.UpdateWallet(walletUpdateDTO)

	assert.Error(t, err)
	assert.Nil(t, updatedWallet)
	assert.Equal(t, serviceerror.WalletDoesntExist, err)
}

func TestDeleteWallet_Success(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWalletRepository := mock.NewMockWalletRepository(ctl)
	mockManager := &repository.Manager{
		Wallet: mockWalletRepository,
	}

	walletService := wallet.NewWalletService(mockManager)

	walletDeleteDTO := &model.WalletDeleteDTO{
		Id: 1,
	}

	mockWalletRepository.
		EXPECT().
		WalletBelongsToUser(walletDeleteDTO.Id, walletDeleteDTO.UserId).
		Times(1).
		Return(true)

	mockWalletRepository.
		EXPECT().
		DeleteWallet(walletDeleteDTO.Id).
		Times(1).
		Return(nil)

	err := walletService.DeleteWallet(walletDeleteDTO)

	assert.NoError(t, err)
}

func TestDeleteWallet_WalletDoesntBelongToUser_Error(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWalletRepository := mock.NewMockWalletRepository(ctl)
	mockManager := &repository.Manager{
		Wallet: mockWalletRepository,
	}

	walletService := wallet.NewWalletService(mockManager)

	walletDeleteDTO := &model.WalletDeleteDTO{
		Id:     1,
		UserId: 1,
	}

	mockWalletRepository.
		EXPECT().
		WalletBelongsToUser(walletDeleteDTO.Id, walletDeleteDTO.UserId).
		Times(1).
		Return(false)

	err := walletService.DeleteWallet(walletDeleteDTO)

	assert.Error(t, err)
	assert.Equal(t, serviceerror.WalletDoesntBelongToUser, err)
}
