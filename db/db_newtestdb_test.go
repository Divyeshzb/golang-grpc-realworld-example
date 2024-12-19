package db

import (
	"os"
	"sync"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestNewTestDB(t *testing.T) {

	tests := []struct {
		name          string
		setupFunc     func()
		cleanupFunc   func()
		expectedError bool
		validateFunc  func(*testing.T, *gorm.DB, error)
	}{
		{
			name: "Successful Database Connection",
			setupFunc: func() {

				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "test_user")
				os.Setenv("DB_PASSWORD", "test_password")
				os.Setenv("DB_NAME", "test_db")
				os.Setenv("DB_PORT", "3306")
			},
			cleanupFunc: func() {
				os.Clearenv()
			},
			expectedError: false,
			validateFunc: func(t *testing.T, db *gorm.DB, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, db)

				maxIdle := db.DB().MaxIdleConns
				assert.Equal(t, 3, maxIdle)
			},
		},
		{
			name: "Missing Environment File",
			setupFunc: func() {

			},
			cleanupFunc: func() {
				os.Clearenv()
			},
			expectedError: true,
			validateFunc: func(t *testing.T, db *gorm.DB, err error) {
				assert.Error(t, err)
				assert.Nil(t, db)
			},
		},
		{
			name: "Invalid Database Credentials",
			setupFunc: func() {
				os.Setenv("DB_HOST", "invalid_host")
				os.Setenv("DB_USER", "invalid_user")
				os.Setenv("DB_PASSWORD", "invalid_password")
				os.Setenv("DB_NAME", "invalid_db")
				os.Setenv("DB_PORT", "3306")
			},
			cleanupFunc: func() {
				os.Clearenv()
			},
			expectedError: true,
			validateFunc: func(t *testing.T, db *gorm.DB, err error) {
				assert.Error(t, err)
				assert.Nil(t, db)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			defer func() {
				if tt.cleanupFunc != nil {
					tt.cleanupFunc()
				}
			}()

			db, err := NewTestDB()

			tt.validateFunc(t, db, err)

			if db != nil {
				db.Close()
			}
		})
	}

	t.Run("Concurrent Access", func(t *testing.T) {

		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_USER", "test_user")
		os.Setenv("DB_PASSWORD", "test_password")
		os.Setenv("DB_NAME", "test_db")
		os.Setenv("DB_PORT", "3306")

		defer os.Clearenv()

		var wg sync.WaitGroup
		numGoroutines := 5
		results := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				db, err := NewTestDB()
				if err != nil {
					results <- err
					return
				}
				defer db.Close()
				results <- nil
			}()
		}

		wg.Wait()
		close(results)

		for err := range results {
			assert.NoError(t, err)
		}
	})

	t.Run("Multiple Sequential Connections", func(t *testing.T) {

		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_USER", "test_user")
		os.Setenv("DB_PASSWORD", "test_password")
		os.Setenv("DB_NAME", "test_db")
		os.Setenv("DB_PORT", "3306")

		defer os.Clearenv()

		for i := 0; i < 3; i++ {
			db, err := NewTestDB()
			assert.NoError(t, err)
			assert.NotNil(t, db)
			if db != nil {
				db.Close()
			}
		}
	})
}
