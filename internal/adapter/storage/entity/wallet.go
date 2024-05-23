package entity

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type Wallet struct {
	Id            uint64          `json:"id" gorm:"primarykey"`
	UserId        uint64          `json:"userId" gorm:"not null;uniqueIndex:idx_userid_name" validate:"required"`
	Name          string          `json:"name" gorm:"not null;uniqueIndex:idx_userid_name" validate:"required,min=3,max=128"`
	Description   *string         `json:"description" gorm:"null" validate:"max=256"`
	Currency      string          `json:"currency" gorm:"not null" validate:"required,len=3"`
	InitialAmount decimal.Decimal `json:"initialAmount" gorm:"not null" validate:"required"`
	CreatedAt     time.Time       `json:"createdAt" gorm:"<-:create"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt  `json:"-" gorm:"index"`
}

func (Wallet) TableName() string { return "portmonetka.wallets" }
