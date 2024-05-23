package storage

import "github.com/khivuksergey/portmonetka.wallet/internal/core/port/repository"

type IDB interface {
	InitRepositoryManager() *repository.Manager
	Close() error
}
