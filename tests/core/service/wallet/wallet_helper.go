package wallet

import "github.com/shopspring/decimal"

func stringPtr(s string) *string {
	return &s
}

func decimalPtr(d decimal.Decimal) *decimal.Decimal {
	return &d
}
