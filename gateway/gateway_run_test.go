package main

import (
	"flag"
	"net/http"
	"testing"
	"time"
)

func TestRun(t *testing.T) {

	tests := []struct {
		name           string
		endpoint       string
		setupMock      func()
		expectedError  bool
		errorContains  string
		contextTimeout time.Duration
	}{
		{
			name:           "Successful Gateway Server Initialization",
			endpoint:       "localhost:50051",
			contextTimeout: 2 * time.Second,
			setupMock:      func() {},
			expectedError:  false,
		},
		{
			name:           "Invalid Endpoint Configuration",
			endpoint:       "invalid:port",
			contextTimeout: 1 * time.Second,
			setupMock:      func() {},
			expectedError:  true,
			errorContains:  "connection refused",
		},
		{
			name:           "Context Cancellation",
			endpoint:       "localhost:50051",
			contextTimeout: 100 * time.Millisecond,
			setupMock:      func() {},
			expectedError:  true,
			errorContains:  "context canceled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Log("Setting up test:", tt.name)
			flag.Set("endpoint", tt.endpoint)
			tt.setupMock()

			errChan := make(chan error, 1)

			go func() {
				err := run()
				errChan <- err
			}()

			select {
			case err := <-errChan:
				if tt.expectedError {
					if err == nil {
						t.Errorf("Expected error but got nil")
					} else if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
						t.Errorf("Expected error containing '%s', got '%s'", tt.errorContains, err.Error())
					}
				} else if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			case <-time.After(tt.contextTimeout):
				if !tt.expectedError {

					t.Log("Server started successfully")

					client := &http.Client{Timeout: 1 * time.Second}
					resp, err := client.Get("http://localhost:3000/health")
					if err != nil {
						t.Logf("Server health check failed: %v", err)
					} else {
						resp.Body.Close()
						t.Log("Server health check passed")
					}
				}
			}
		})
	}
}
func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[len(s)-len(substr):] == substr
}
