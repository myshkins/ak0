package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckRateLimit(t *testing.T) {
	crl := NewClientRateLimiters()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	limitedHandler := CheckRateLimit(crl, handler)

	// create a mock server
	server := httptest.NewServer(limitedHandler)
	defer server.Close()

	client := server.Client()

	// test for allowed requests
	for i := 0; i < 10; i++ {
		req, _ := http.NewRequest("GET", server.URL, nil)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}
	}

	// Test blocking requests over the TestCheckRateLimit
	req, _ := http.NewRequest("GET", server.URL, nil)
	req.RemoteAddr = "192.168.1.1"
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status Too Many Requests, got %v", resp.Status)
	}
}

func TestGetLimiter(t *testing.T) {
	crl := NewClientRateLimiters()

	ip := "192.168.1.1"
	limiter := getLimiter(crl, ip)

	if limiter == nil {
		t.Fatalf("Expected limiter to be created, got nil")
	}

	limiter2 := getLimiter(crl, ip)
	if limiter != limiter2 {
		t.Fatal("Expected the same limiter instance")
	}
}

// todo: fix this. right now it only runs after 1min, so test doesn't work
// func TestCleanupRateLimiters(t *testing.T) {
// 	crl := NewClientRateLimiters()

// 	// create a stale limiter
// 	crl.ClientLimiters["192.168.1.1"] = &ClientRateLimiter{
// 		limiter:  rate.NewLimiter(1, 1),
// 		lastSeen: time.Now().Add(-4 * time.Minute),
// 	}

// 	go CleanupRateLimiters(context.Background(), crl)

// 	time.Sleep(2 * time.Second)

// 	crl.Mu.Lock()
// 	defer crl.Mu.Unlock()
// 	if _, exists := crl.ClientLimiters["192.168.1.1"]; exists {
// 		t.Fatal("Expected old limiter to be cleaned up")
// 	}
// }
