package model

import (
	"github.com/go-playground/validator/v10"
)

func GetWalletValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	//TODO Add custom validation messages
	//v.RegisterStructValidation(validateWalletUpdate, WalletUpdateDTO{})

	//// Add custom message for the validation error
	//v.RegisterTranslation("atleastonefieldrequired", en.Translations, func(ut validator.Translator) error {
	//	return ut.Add("atleastonefieldrequired", "At least one field must be provided", true)
	//}, func(ut validator.Translator, fe validator.FieldError) string {
	//	t, _ := ut.T("atleastonefieldrequired", fe.Field())
	//	return t
	//})

	return v
}

//func validateWalletUpdate(sl validator.StructLevel) {
//	wallet := sl.Current().Interface().(WalletUpdateDTO)
//
//	if wallet.Name == nil && wallet.Description == nil && wallet.Currency == nil && wallet.InitialAmount == nil {
//		sl.ReportError(wallet, "WalletUpdateDTO", "", "atleastonefieldrequired", "")
//	}
//}
