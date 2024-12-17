package db

import (
	"testing"
	"sync"
	"github.com/jinzhu/gorm"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestDropTestDB(t *testing.T) {

	type testCase struct {
		name     string
		db       *gorm.DB
		wantErr  bool
		setupFn  func() (*gorm.DB, sqlmock.Sqlmock, error)
		validate func(*testing.T, *gorm.DB, error)
	}

	createMockDB := func() (*gorm.DB, sqlmock.Sqlmock, error) {
		sqlDB, mock, err := sqlmock.New()
		if err != nil {
			return nil, nil, err
		}
		gormDB, err := gorm.Open("mysql", sqlDB)
		if err != nil {
			return nil, nil, err
		}
		return gormDB, mock, nil
	}

	tests := []testCase{
		{
			name: "Successful Database Closure",
			setupFn: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				return createMockDB()
			},
			validate: func(t *testing.T, db *gorm.DB, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}

				sqlDB := db.DB()
				if err := sqlDB.Ping(); err == nil {
					t.Error("Expected database to be closed")
				}
			},
		},
		{
			name: "Nil Database Instance",
			setupFn: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				return nil, nil, nil
			},
			validate: func(t *testing.T, db *gorm.DB, err error) {
				if err != nil {
					t.Logf("Expected behavior with nil DB: %v", err)
				}
			},
		},
		{
			name: "Concurrent Database Closure",
			setupFn: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				return createMockDB()
			},
			validate: func(t *testing.T, db *gorm.DB, err error) {
				var wg sync.WaitGroup
				for i := 0; i < 5; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						err := DropTestDB(db)
						if err != nil {
							t.Errorf("Concurrent closure failed: %v", err)
						}
					}()
				}
				wg.Wait()
			},
		},
		{
			name: "Database with Active Transaction",
			setupFn: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := createMockDB()
				if err != nil {
					return nil, nil, err
				}
				mock.ExpectBegin()
				return db.Begin(), mock, nil
			},
			validate: func(t *testing.T, db *gorm.DB, err error) {
				if err != nil {
					t.Errorf("Failed to handle active transaction: %v", err)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, _, err := tc.setupFn()
			if err != nil && !tc.wantErr {
				t.Fatalf("Setup failed: %v", err)
			}

			err = DropTestDB(db)

			tc.validate(t, db, err)
		})
	}
}
