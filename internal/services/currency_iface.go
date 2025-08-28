package services

import "context"

type CurrencyConverter interface {
	Convert(ctx context.Context, amount float64, from, to string) (float64, float64, error)
}
