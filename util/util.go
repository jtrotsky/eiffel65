package util

import "github.com/google/uuid"

// GenerateUUID creates a new random unique ID.
func GenerateUUID() string {
	id, _ := uuid.NewRandom()
	return id.String()
}

// Steam Currencies
// "USD": 1,  # United States dollar
// "GBP": 2,  # British pound sterling
// "EUR": 3,  # The euro
// "CHF": 4,  # Swiss franc
// "RUB": 5,  # Russian ruble
// "PLN": 6,  # Polish złoty
// "BRL": 7,  # Brazilian real
// "JPY": 8,  # Japanese yen
// "SEK": 9,  # Swedish krona
// "IDR": 10,  # Indonesian rupiah
// "MYR": 11,  # Malaysian ringgit
// "BWP": 12,  # Botswana pula
// "SGD": 13,  # Singapore dollar
// "THB": 14,  # Thai baht
// "VND": 15,  # Vietnamese dong
// "KRW": 16,  # South Korean won
// "TRY": 17,  # Turkish lira
// "UAH": 18,  # Ukrainian hryvnia
// "MXN": 19,  # Mexican Peso
// "CAD": 20,  # Canadian dollar
// "AUD": 21,  # Australian dollar
// "NZD": 22,  # New Zealand dollar
// "CNY": 23,  # Chinese yuan
// "INR": 24,  # Indian rupee
// "CLP": 25,  # Chilean peso
// "PEN": 26,  # Peruvian sol
// "COP": 27,  # Colombian peso
// "ZAR": 28,  # South African rand
// "HKD": 29,  # Hong Kong dollar
// "TWD": 30,  # New Taiwan dollar
// "SAR": 31,  # Saudi riyal
// "AED": 32  # United Arab Emirates dirham
