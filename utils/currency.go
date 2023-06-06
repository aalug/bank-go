package utils

// all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	PLN = "PLN"
)

// IsSupportedCurrency checks if currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, PLN:
		return true
	}
	return false
}

