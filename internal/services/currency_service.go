package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type CurrencyService struct {
	redis  *redis.Client
	apiURL string
	ttl    time.Duration
}

func NewCurrencyService(r *redis.Client, apiURL string, ttl time.Duration) *CurrencyService {
	return &CurrencyService{redis: r, apiURL: apiURL, ttl: ttl}
}

func cacheKey(from, to string) string {
	return fmt.Sprintf("fx:%s:%s", strings.ToUpper(from), strings.ToUpper(to))
}

func (s *CurrencyService) Convert(ctx context.Context, amount float64, from, to string) (float64, float64, error) {

	from = strings.ToUpper(from)
	to = strings.ToUpper(to)
	key := cacheKey(from, to)

	if s.redis != nil {
		if val, err := s.redis.Get(ctx, key).Result(); err == nil {
			var rate float64
			if _ = json.Unmarshal([]byte(val), &rate); rate != 0 {
				return amount * rate, rate, nil
			}
		} else if err != redis.Nil {
			log.Printf("Redis error: %v", err)
		}
	}
	url := fmt.Sprintf("%s/latest/%s", s.apiURL, from)
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("failed to fetch exchange rate: %s", resp.Status)
	}

	var data struct {
		ConversionRates map[string]float64 `json:"conversion_rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, err
	}

	rate, ok := data.ConversionRates[to]

	if !ok {
		return 0, 0, fmt.Errorf("unsupported currency: %s", to)
	}

	if s.redis != nil {
		if val, err := json.Marshal(rate); err == nil {
			if err := s.redis.Set(ctx, key, val, s.ttl).Err(); err != nil {
				log.Printf("Redis error: %v", err)
			}
		}
	}

	return amount * rate, rate, nil
}
