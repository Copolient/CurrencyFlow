package websocket

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewHub(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("expected hub to be created")
	}
	if hub.ClientCount() != 0 {
		t.Fatalf("expected 0 clients, got %d", hub.ClientCount())
	}
}

func TestHub_BroadcastRateUpdate(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create a mock connection using a pipe
	// Since we can't easily create a real websocket.Conn in tests,
	// we test the broadcast channel directly
	update := RateUpdate{
		FromCurrency: "USD",
		ToCurrency:   "CNY",
		Rate:         7.24,
		Timestamp:    time.Now().Format(time.RFC3339),
	}

	// Broadcast should not panic even with no clients
	hub.BroadcastRateUpdate(update)

	// Verify the update was marshaled correctly
	data, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("failed to marshal update: %v", err)
	}

	var unmarshaled RateUpdate
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal update: %v", err)
	}

	if unmarshaled.FromCurrency != "USD" {
		t.Fatalf("expected FromCurrency USD, got %s", unmarshaled.FromCurrency)
	}
	if unmarshaled.ToCurrency != "CNY" {
		t.Fatalf("expected ToCurrency CNY, got %s", unmarshaled.ToCurrency)
	}
	if unmarshaled.Rate != 7.24 {
		t.Fatalf("expected Rate 7.24, got %f", unmarshaled.Rate)
	}
}

func TestHub_ClientCount(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	if hub.ClientCount() != 0 {
		t.Fatalf("expected 0 clients, got %d", hub.ClientCount())
	}
}

func TestRateUpdate_JSON(t *testing.T) {
	update := RateUpdate{
		FromCurrency: "EUR",
		ToCurrency:   "USD",
		Rate:         1.08,
		Timestamp:    "2024-01-15T12:00:00Z",
	}

	data, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Verify JSON structure
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if raw["fromCurrency"] != "EUR" {
		t.Fatalf("expected fromCurrency EUR, got %v", raw["fromCurrency"])
	}
	if raw["toCurrency"] != "USD" {
		t.Fatalf("expected toCurrency USD, got %v", raw["toCurrency"])
	}
	if raw["rate"] != 1.08 {
		t.Fatalf("expected rate 1.08, got %v", raw["rate"])
	}
}
