package db

import (
	"database/sql"
	"sync"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestAutoMigrate(t *testing.T) {

	type testCase struct {
		name        string
		setupDB     func() (*gorm.DB, sqlmock.Sqlmock, error)
		expectError bool
		concurrent  bool
	}

	tests := []testCase{
		{
			name: "Successful Migration",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}

				mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("ALTER TABLE").WillReturnResult(sqlmock.NewResult(1, 1))

				gormDB, err := gorm.Open("mysql", db)
				return gormDB, mock, err
			},
			expectError: false,
			concurrent:  false,
		},
		{
			name: "Database Connection Error",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}

				mock.ExpectExec("CREATE TABLE").WillReturnError(sql.ErrConnDone)

				gormDB, err := gorm.Open("mysql", db)
				return gormDB, mock, err
			},
			expectError: true,
			concurrent:  false,
		},
		{
			name: "Concurrent Migration",
			setupDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}

				mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("ALTER TABLE").WillReturnResult(sqlmock.NewResult(1, 1))

				gormDB, err := gorm.Open("mysql", db)
				return gormDB, mock, err
			},
			expectError: false,
			concurrent:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := tc.setupDB()
			if err != nil {
				t.Fatalf("Failed to setup test database: %v", err)
			}
			defer db.Close()

			if tc.concurrent {

				var wg sync.WaitGroup
				errChan := make(chan error, 3)

				for i := 0; i < 3; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						err := AutoMigrate(db)
						errChan <- err
					}()
				}

				wg.Wait()
				close(errChan)

				for err := range errChan {
					if tc.expectError {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
					}
				}
			} else {

				err = AutoMigrate(db)

				if tc.expectError {
					assert.Error(t, err)
					t.Logf("Expected error occurred: %v", err)
				} else {
					assert.NoError(t, err)
					t.Log("Migration completed successfully")
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
