package util


const (
	USD = "USD"
	EUR = "EUR"
	CNY = "CNY"
)

func IsValidCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CNY:
		return true
	default:
		return false
	}
}
