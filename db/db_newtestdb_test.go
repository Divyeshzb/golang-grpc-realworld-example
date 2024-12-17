package db

import (
	"os"
	"sync"
	"testing"
	"time"
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

	backupEnvFile := func() error {
		content, err := os.ReadFile("../env/test.env")
		if err != nil {
			return err
		}
		return os.WriteFile("../env/test.env.backup", content, 0644)
	}

	restoreEnvFile := func() error {
		content, err := os.ReadFile("../env/test.env.backup")
		if err != nil {
			return err
		}
		err = os.WriteFile("../env/test.env", content, 0644)
		if err != nil {
			return err
		}
		return os.Remove("../env/test.env.backup")
	}

	tests := []testCase{
		{
			name: "Successful Database Connection",
			setupFunc: func() {

				_ = backupEnvFile()
			},
			cleanupFunc: func() {

				_ = restoreEnvFile()
			},
			expectedError: false,
			validateFunc: func(t *testing.T, db *gorm.DB, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, db)

				sqlDB := db.DB()
				maxIdle := sqlDB.Stats().MaxOpenConnections
				assert.Equal(t, 3, maxIdle)

				err = sqlDB.Ping()
				assert.NoError(t, err)
			},
		},
		{
			name: "Missing Environment File",
			setupFunc: func() {

				_ = os.Rename("../env/test.env", "../env/test.env.tmp")
			},
			cleanupFunc: func() {

				_ = os.Rename("../env/test.env.tmp", "../env/test.env")
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
				_ = backupEnvFile()

				invalidEnv := []byte("DB_HOST=invalid\nDB_USER=invalid\nDB_PASSWORD=invalid\nDB_NAME=invalid\nDB_PORT=3306")
				_ = os.WriteFile("../env/test.env", invalidEnv, 0644)
			},
			cleanupFunc: func() {
				_ = restoreEnvFile()
			},
			expectedError: true,
			validateFunc: func(t *testing.T, db *gorm.DB, err error) {
				assert.Error(t, err)
				assert.Nil(t, db)
			},
		},
		{
			name: "Concurrent Access",
			setupFunc: func() {
				_ = backupEnvFile()
			},
			cleanupFunc: func() {
				_ = restoreEnvFile()
			},
			expectedError: false,
			validateFunc: func(t *testing.T, _ *gorm.DB, _ error) {
				var wg sync.WaitGroup
				concurrentAccess := 5
				errorChan := make(chan error, concurrentAccess)

				for i := 0; i < concurrentAccess; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						db, err := NewTestDB()
						if err != nil {
							errorChan <- err
							return
						}
						if db != nil {
							defer db.Close()
						}
					}()
				}

				wg.Wait()
				close(errorChan)

				for err := range errorChan {
					assert.NoError(t, err)
				}
			},
		},
		{
			name: "Multiple Sequential Connections",
			setupFunc: func() {
				_ = backupEnvFile()
			},
			cleanupFunc: func() {
				_ = restoreEnvFile()
			},
			expectedError: false,
			validateFunc: func(t *testing.T, _ *gorm.DB, _ error) {

				for i := 0; i < 3; i++ {
					db, err := NewTestDB()
					assert.NoError(t, err)
					assert.NotNil(t, db)
					if db != nil {
						defer db.Close()
					}

					time.Sleep(100 * time.Millisecond)
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
