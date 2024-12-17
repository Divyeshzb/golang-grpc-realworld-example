package db

import (
	"os"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	type testCase struct {
		name    string
		envVars map[string]string
		wantErr bool
		errMsg  string
	}

	tests := []testCase{
		{
			name: "Successful Database Connection",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			wantErr: false,
		},
		{
			name: "Missing DB_HOST",
			envVars: map[string]string{
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			wantErr: true,
			errMsg:  "$DB_HOST is not set",
		},
		{
			name: "Missing DB_USER",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			wantErr: true,
			errMsg:  "$DB_USER is not set",
		},
		{
			name: "Database Server Unavailable",
			envVars: map[string]string{
				"DB_HOST":     "nonexistent",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
				"DB_PORT":     "3306",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			os.Clearenv()

			for k, v := range tc.envVars {
				os.Setenv(k, v)
			}

			start := time.Now()
			db, err := New()
			duration := time.Since(start)

			t.Logf("Test '%s' took %v to execute", tc.name, duration)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.errMsg != "" {
					assert.Contains(t, err.Error(), tc.errMsg)
				}
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)

				if db != nil {
					sqlDB := db.DB()
					maxIdle := sqlDB.Stats().MaxOpenConnections
					assert.Equal(t, 3, maxIdle, "MaxIdleConns should be 3")

					db.Close()
				}
			}
		})
	}
}
