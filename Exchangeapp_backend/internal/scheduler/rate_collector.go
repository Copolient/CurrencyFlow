package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"exchangeapp/internal/service"
	ws "exchangeapp/internal/websocket"

	"github.com/robfig/cron/v3"
)

type RateCollector struct {
	rateHistorySvc *service.RateHistoryService
	hub            *ws.Hub
	httpClient     *http.Client
	cron           *cron.Cron
}

type exchangeRateAPIResponse struct {
	Base  string             `json:"base_code"`
	Rates map[string]float64 `json:"conversion_rates"`
}

func NewRateCollector(rateHistorySvc *service.RateHistoryService, hub *ws.Hub) *RateCollector {
	return &RateCollector{
		rateHistorySvc: rateHistorySvc,
		hub:            hub,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cron: cron.New(),
	}
}

// Currencies to track
var trackedPairs = []struct{ From, To string }{
	{"USD", "CNY"}, {"USD", "EUR"}, {"USD", "GBP"}, {"USD", "JPY"},
	{"EUR", "CNY"}, {"EUR", "GBP"}, {"EUR", "JPY"},
	{"GBP", "CNY"}, {"GBP", "JPY"},
	{"CNY", "JPY"},
}

func (rc *RateCollector) Start() {
	// Run every 30 minutes
	rc.cron.AddFunc("*/30 * * * *", func() {
		rc.collect()
	})

	// Run once on startup
	go rc.collect()

	rc.cron.Start()
	log.Println("Rate collector started (every 30 min)")
}

func (rc *RateCollector) Stop() {
	rc.cron.Stop()
}

func (rc *RateCollector) collect() {
	ctx := context.Background()

	// Fetch USD-based rates from free API
	rates, err := rc.fetchRates("USD")
	if err != nil {
		log.Printf("RateCollector: failed to fetch rates: %v", err)
		return
	}

	// Cross-calculate and record
	for _, pair := range trackedPairs {
		var rate float64
		if pair.From == "USD" {
			r, ok := rates[pair.To]
			if !ok {
				continue
			}
			rate = r
		} else if pair.To == "USD" {
			r, ok := rates[pair.From]
			if !ok || r == 0 {
				continue
			}
			rate = 1.0 / r
		} else {
			// Cross rate: From→USD→To
			fromRate, ok1 := rates[pair.From]
			toRate, ok2 := rates[pair.To]
			if !ok1 || !ok2 || fromRate == 0 {
				continue
			}
			rate = toRate / fromRate
		}

		if err := rc.rateHistorySvc.RecordRate(ctx, pair.From, pair.To, rate); err != nil {
			log.Printf("RateCollector: failed to record %s/%s: %v", pair.From, pair.To, err)
			continue
		}

		// Broadcast via WebSocket
		if rc.hub != nil {
			rc.hub.BroadcastRateUpdate(ws.RateUpdate{
				FromCurrency: pair.From,
				ToCurrency:   pair.To,
				Rate:         rate,
				Timestamp:    time.Now().Format(time.RFC3339),
			})
		}
	}

	log.Println("RateCollector: rates collected successfully")
}

func (rc *RateCollector) fetchRates(base string) (map[string]float64, error) {
	// Using exchangerate-api.com free tier (no key needed for basic)
	url := fmt.Sprintf("https://open.er-api.com/v6/latest/%s", base)

	resp, err := rc.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Limit response body to 1MB
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	var apiResp exchangeRateAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return apiResp.Rates, nil
}
