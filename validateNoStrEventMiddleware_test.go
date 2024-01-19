package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateNoStrEventMiddleware(t *testing.T) {
	tests := []struct {
		name     string
		payload  RequestPayload
		expected int
	}{
		{
			name: "Valid TAHUB_CREATE_USER",
			payload: RequestPayload{
				Kind:    1,
				Content: "TAHUB_CREATE_USER",
			},
			expected: http.StatusOK,
		},
		{
			name: "Valid TAHUB_GET_BALANCES",
			payload: RequestPayload{
				Kind:    1,
				Content: "TAHUB_GET_BALANCES",
			},
			expected: http.StatusOK,
		},
		{
			name: "Valid TAHUB_RECEIVE_ADDRESS_FOR_ASSET",
			payload: RequestPayload{
				Kind:    1,
				Content: "TAHUB_RECEIVE_ADDRESS_FOR_ASSET",
				Ta:      "some_value",
				Amt:     10.0,
			},
			expected: http.StatusOK,
		},
		{
			name: "Valid TAHUB_SEND_ASSET",
			payload: RequestPayload{
				Kind: 1,
				Content: "TAHUB_SEND_ASSET",
				Addr: "some_address",
				Fee:  5.0,
			},
			expected: http.StatusOK,
		},
		{
			name: "Invalid event content",
			payload: RequestPayload{
				Kind:    1,
				Content: "invalid_event",
			},
			expected: http.StatusBadRequest,
		},
	
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := validateNoStrEventMiddleware(processNoStrEvent)

			// Create a mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(handler))
			defer server.Close()

			requestBody, err := json.Marshal(test.payload)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.Post(server.URL, "application/json", strings.NewReader(string(requestBody)))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			
			if resp.StatusCode != test.expected {
				t.Errorf("Expected status code %d, got %d", test.expected, resp.StatusCode)
			}
		})
	}
}
