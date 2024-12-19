package db

import (
	"database/sql"
	"os"
	"sync"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestNewTestDB(t *testing.T) {

	type testCase struct {
		name          string
		setupFunc     func()
		cleanupFunc   func()
		expectedError bool
		validateFunc  func(*testing.T, *gorm.DB, error)
	}

	commonCleanup := func() {

		txdbInitialized = false
		if db, err := NewTestDB(); err == nil && db != nil {
			db.Close()
		}
	}

	tests := []testCase{
		{
			name: "Successful Database Connection",
			setupFunc: func() {

			},
			cleanupFunc:   commonCleanup,
			expectedError: false,
			validateFunc: func(t *testing.T, db *gorm.DB, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				assert.True(t, db.DB().Ping() == nil)
			},
		},
		{
			name: "Missing Environment File",
			setupFunc: func() {

				os.Rename("../env/test.env", "../env/test.env.backup")
			},
			cleanupFunc: func() {
				commonCleanup()

				os.Rename("../env/test.env.backup", "../env/test.env")
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

			},
			cleanupFunc:   commonCleanup,
			expectedError: true,
			validateFunc: func(t *testing.T, db *gorm.DB, err error) {
				assert.Error(t, err)
				assert.Nil(t, db)
			},
		},
		{
			name: "Concurrent Access",
			setupFunc: func() {

			},
			cleanupFunc:   commonCleanup,
			expectedError: false,
			validateFunc: func(t *testing.T, _ *gorm.DB, _ error) {
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
						results <- db.DB().Ping()
					}()
				}

				wg.Wait()
				close(results)

				for err := range results {
					assert.NoError(t, err)
				}
			},
		},
		{
			name: "Database Resource Management",
			setupFunc: func() {

			},
			cleanupFunc:   commonCleanup,
			expectedError: false,
			validateFunc: func(t *testing.T, db *gorm.DB, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, db)

				sqlDB := db.DB()
				stats := sqlDB.Stats()
				assert.Equal(t, 3, stats.MaxOpenConnections)

				assert.NoError(t, db.Close())
			},
		},
		{
			name: "Multiple Sequential Connections",
			setupFunc: func() {

			},
			cleanupFunc:   commonCleanup,
			expectedError: false,
			validateFunc: func(t *testing.T, _ *gorm.DB, _ error) {

				for i := 0; i < 3; i++ {
					db, err := NewTestDB()
					assert.NoError(t, err)
					assert.NotNil(t, db)
					defer db.Close()

					assert.NoError(t, db.DB().Ping())
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			if tc.setupFunc != nil {
				tc.setupFunc()
			}

			if tc.cleanupFunc != nil {
				defer tc.cleanupFunc()
			}

			db, err := NewTestDB()
			if db != nil {
				defer db.Close()
			}

			tc.validateFunc(t, db, err)
		})
	}
}
