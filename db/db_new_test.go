package db

import (
	"database/sql"
	"os"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	type testCase struct {
		name     string
		envVars  map[string]string
		mockFunc func()
		wantErr  bool
		errMsg   string
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
			name: "Invalid Credentials",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_USER":     "invalid",
				"DB_PASSWORD": "invalid",
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

			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			db, err := New()

			if tc.wantErr {
				assert.Error(t, err)
				if tc.errMsg != "" {
					assert.Contains(t, err.Error(), tc.errMsg)
				}
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)

				sqlDB := db.DB()
				maxIdle := sqlDB.Stats().MaxOpenConnections
				assert.Equal(t, 3, maxIdle)

				db.Close()
			}
		})
	}
}
func TestNewConcurrent(t *testing.T) {

	envVars := map[string]string{
		"DB_HOST":     "localhost",
		"DB_USER":     "testuser",
		"DB_PASSWORD": "testpass",
		"DB_NAME":     "testdb",
		"DB_PORT":     "3306",
	}

	for k, v := range envVars {
		os.Setenv(k, v)
	}

	numGoroutines := 5
	done := make(chan bool)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			db, err := New()
			if err == nil {
				db.Close()
			}
			done <- true
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}
