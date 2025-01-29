package middleware

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"reflect"
// 	"testing"
// )

// var expected = []{}

// func TestFilterBots(t *testing.T) {
// 	endpoints := parseEndpointConfig([]byte(configYaml))

// 	if !reflect.DeepEqual(expected_endpoints, endpoints) {
// 		t.Errorf("\nendpoints = \n%+v, \nwant \n%+v", endpoints, expected_endpoints)
// 	}
// }

// func TestPing(t *testing.T) {
// 	// create a mock server
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 	}))
// 	defer server.Close()

// 	endpoints := []endpoint{
// 		{Name: "MockEndpoint", Url: server.URL},
// 	}

// 	hc := &HealthCheckClient{
// 		httpClient:   server.Client(),
// 		endpoints:    endpoints,
// 		timeInterval: 1,
// 		stats:        make(map[string]map[string]int),
// 	}

// 	for _, e := range hc.endpoints {
// 		hc.stats[e.Name] = make(map[string]int)
// 	}

// 	err := hc.ping()
// 	if err != nil {
// 		t.Fatalf("ping returned an error: %v", err)
// 	}

// 	if hc.stats["MockEndpoint"]["up"] != 1 {
// 		t.Errorf("expected 'up' count to be 1, got %d", hc.stats["MockEndpoint"]["up"])
// 	}
// }
