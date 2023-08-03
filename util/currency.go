package util

// check if a currency is supported or not

// constants for all supported currency
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

// isSupportedCurrency return true if currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	//else
	return false
}
