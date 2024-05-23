package error

import (
	"errors"
	"fmt"
)

var (
	WalletAlreadyExists          = errors.New("wallet with this name already exists")
	WalletDoesntExist            = errors.New("wallet with this id doesn't exists")
	WalletDoesntBelongToUser     = errors.New("wallet with this id doesn't belong to user")
	AtLeastOneFieldIsRequired    = errors.New("at least one field for updating wallet is required")
	WalletNameLengthError        = errors.New("wallet name must be from 3 to 128 symbols long")
	WalletDescriptionLengthError = errors.New("wallet description must be less than 256 symbols long")
	WalletCurrencyError          = errors.New("wallet currency must be 3 symbols long")
)

const (
	InvalidInputData   = "invalid input data"
	CannotCreateWallet = "cannot create wallet"
	CannotGetWallets   = "cannot retrieve wallets"
	CannotUpdateWallet = "cannot update wallet"
	CannotDeleteWallet = "cannot delete wallet"
)

type ErrorMessage string

func (m *ErrorMessage) Append(errMessage string) {
	if *m != "" {
		*m += "; "
	}
	*m += ErrorMessage(errMessage)
}

func (m *ErrorMessage) ToError() error {
	if *m == "" {
		return nil
	}
	return fmt.Errorf(fmt.Sprint(*m))
}
