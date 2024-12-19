package db

import (
	"sync"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func TestDropTestDB(t *testing.T) {

	type testCase struct {
		name          string
		setupDB       func() (*gorm.DB, sqlmock.Sqlmock, error)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successfully Close Database Connection",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				gormDB, err := gorm.Open("mysql", db)
				if err != nil {
					return nil, nil, err
				}
				mock.ExpectClose()
				return gormDB, mock, nil
			},
			expectedError: nil,
		},
		{
			name: "Handle Nil Database Instance",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				return nil, nil, nil
			},
			expectedError: nil,
		},
		{
			name: "Handle Already Closed Database",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				gormDB, err := gorm.Open("mysql", db)
				if err != nil {
					return nil, nil, err
				}
				gormDB.Close()
				mock.ExpectClose()
				return gormDB, mock, nil
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := tc.setupDB()
			if err != nil {
				t.Fatalf("Error setting up test database: %v", err)
			}

			err = DropTestDB(db)

			if err != tc.expectedError {
				t.Errorf("Expected error %v, got %v", tc.expectedError, err)
			}

			if mock != nil {
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("There were unfulfilled expectations: %s", err)
				}
			}
		})
	}

	t.Run("Concurrent Database Closure", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock database: %v", err)
		}
		gormDB, err := gorm.Open("mysql", db)
		if err != nil {
			t.Fatalf("Error creating GORM database: %v", err)
		}

		mock.ExpectClose()

		var wg sync.WaitGroup
		concurrentCalls := 5

		errorChan := make(chan error, concurrentCalls)

		for i := 0; i < concurrentCalls; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := DropTestDB(gormDB)
				errorChan <- err
			}()
		}

		wg.Wait()
		close(errorChan)

		for err := range errorChan {
			if err != nil {
				t.Errorf("Concurrent closure resulted in error: %v", err)
			}
		}
	})

	t.Run("Database with Active Transactions", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock database: %v", err)
		}
		gormDB, err := gorm.Open("mysql", db)
		if err != nil {
			t.Fatalf("Error creating GORM database: %v", err)
		}

		mock.ExpectBegin()
		mock.ExpectClose()

		tx := gormDB.Begin()
		if tx.Error != nil {
			t.Fatalf("Error starting transaction: %v", tx.Error)
		}

		err = DropTestDB(gormDB)
		if err != nil {
			t.Errorf("Expected no error when closing DB with active transaction, got: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	})
}
