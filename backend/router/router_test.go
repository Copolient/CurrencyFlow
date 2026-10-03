package router

import "testing"

func TestSetupRouterRegistersExpectedRoutes(t *testing.T) {
	engine, limiter, authLimiter := SetupRouter(Handlers{}, nil, nil, nil)
	defer limiter.Stop()
	defer authLimiter.Stop()

	registered := make(map[string]bool)
	for _, route := range engine.Routes() {
		registered[route.Method+" "+route.Path] = true
	}

	expected := []string{
		"POST /api/v1/auth/login",
		"POST /api/v1/auth/register",
		"GET /api/v1/exchangeRates",
		"GET /api/v1/rates/history",
		"GET /api/v1/rates/latest",
		"GET /api/v1/posts",
		"GET /api/v1/users/:id",
		"POST /api/v1/ai/analyze",
		"GET /api/v1/ws",
		"GET /healthz",
		"GET /metrics",
		"POST /api/v1/articles",
		"GET /api/v1/articles",
		"GET /api/v1/articles/:id",
		"POST /api/v1/alerts",
		"GET /api/v1/alerts",
		"DELETE /api/v1/alerts/:id",
		"GET /api/v1/notifications",
		"PUT /api/v1/notifications/:id/read",
		"PUT /api/v1/notifications/read-all",
		"GET /api/v1/notifications/unread-count",
		"POST /api/v1/posts",
		"POST /api/v1/posts/:id/like",
		"POST /api/v1/users/:id/follow",
		"DELETE /api/v1/users/:id/follow",
		"GET /api/v1/users/:id/following",
		"POST /api/v1/favorites",
		"GET /api/v1/favorites",
		"DELETE /api/v1/favorites",
		"GET /api/v1/favorites/check",
		"PUT /api/v1/users/profile",
		"POST /api/auth/login",
		"POST /api/auth/register",
		"GET /api/exchangeRates",
	}

	for _, key := range expected {
		if registered[key] {
			continue
		}
		t.Fatalf("expected route to be registered: %s", key)
	}
}
