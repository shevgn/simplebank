package util

const (
	USD = "USD"
	EUR = "EUR"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR:
		return true
	default:
		return false
	}
}
