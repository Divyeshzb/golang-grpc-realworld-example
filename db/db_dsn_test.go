package db

import (
	"os"
	"testing"
)

func Testdsn(t *testing.T) {

	type testCase struct {
		name     string
		envVars  map[string]string
		expected string
		wantErr  string
	}

	tests := []testCase{
		{
			name: "Successfully Generate DSN String",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			expected: "testuser:testpass@(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
			wantErr:  "",
		},
		{
			name: "Missing DB_HOST",
			envVars: map[string]string{
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			expected: "",
			wantErr:  "$DB_HOST is not set",
		},
		{
			name: "Missing DB_USER",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			expected: "",
			wantErr:  "$DB_USER is not set",
		},
		{
			name: "Missing DB_PASSWORD",
			envVars: map[string]string{
				"DB_HOST": "localhost",
				"DB_USER": "testuser",
				"DB_NAME": "testdb",
				"DB_PORT": "3306",
			},
			expected: "",
			wantErr:  "$DB_PASSWORD is not set",
		},
		{
			name: "Missing DB_NAME",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_PORT":     "3306",
			},
			expected: "",
			wantErr:  "$DB_NAME is not set",
		},
		{
			name: "Missing DB_PORT",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			expected: "",
			wantErr:  "$DB_PORT is not set",
		},
		{
			name: "Special Characters in Credentials",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "user@special",
				"DB_PASSWORD": "p@ssw#rd!",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			expected: "user@special:p@ssw#rd!@(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
			wantErr:  "",
		},
		{
			name: "Empty String in DB_HOST",
			envVars: map[string]string{
				"DB_HOST":     "",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			expected: "",
			wantErr:  "$DB_HOST is not set",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			os.Clearenv()

			for key, value := range tc.envVars {
				os.Setenv(key, value)
			}

			got, err := dsn()

			if tc.wantErr != "" {
				if err == nil {
					t.Errorf("dsn() error = nil, wantErr %v", tc.wantErr)
					return
				}
				if err.Error() != tc.wantErr {
					t.Errorf("dsn() error = %v, wantErr %v", err, tc.wantErr)
					return
				}
				t.Logf("Successfully caught expected error: %v", err)
			} else {
				if err != nil {
					t.Errorf("dsn() unexpected error = %v", err)
					return
				}
				if got != tc.expected {
					t.Errorf("dsn() = %v, want %v", got, tc.expected)
					return
				}
				t.Logf("Successfully generated DSN string: %v", got)
			}
		})
	}
}
