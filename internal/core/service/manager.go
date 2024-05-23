package service

import (
	"github.com/khivuksergey/portmonetka.wallet/internal/core/port/repository"
	"github.com/khivuksergey/portmonetka.wallet/internal/core/port/service"
	"github.com/khivuksergey/portmonetka.wallet/internal/core/service/wallet"
)

func NewServiceManager(repositoryManager *repository.Manager) *service.Manager {
	return &service.Manager{
		Wallet: wallet.NewWalletService(repositoryManager),
	}
}
