package utils

import (
	"fmt"
)

type CurrencyConverter interface {
	Convert(from, to string, amount float64) (float64, error)
}

// FixedRateConverter pakai map rate statis (testing/demo)
type FixedRateConverter struct {
	rates map[string]float64
}

func NewFixedRateConverter() *FixedRateConverter {
	return &FixedRateConverter{
		rates: map[string]float64{
			"USD:IDR": 16000,
			"IDR:USD": 1.0 / 16000,
			"USD:EUR": 0.9,
			"EUR:USD": 1.111,
			"USD:JPY": 147.49,
			"JPY:USD": 1.0 / 147.49,
			"EUR:JPY": 165.0,
			"JPY:EUR": 1.0 / 165.0,
		},
	}
}

func (c *FixedRateConverter) Convert(from, to string, amount float64) (float64, error) {
	if from == to {
		return amount, nil
	}
	key := from + ":" + to
	rate, ok := c.rates[key]
	if !ok {
		return 0, fmt.Errorf("unsupported currency conversion %s → %s", from, to)
	}

	result := amount * rate
	return result, nil
}
