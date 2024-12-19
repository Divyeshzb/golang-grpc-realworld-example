package db

import (
	"errors"
	"sync"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"testing"
)

func TestAutoMigrate(t *testing.T) {
	type testCase struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		concurrent  bool
		wantErr     bool
		expectedErr error
	}

	tests := []testCase{
		{
			name: "Successful Migration",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("ALTER TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			concurrent:  false,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Database Connection Error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("CREATE TABLE").WillReturnError(errors.New("connection refused"))
			},
			concurrent:  false,
			wantErr:     true,
			expectedErr: errors.New("connection refused"),
		},
		{
			name: "Concurrent Migration",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("ALTER TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			concurrent:  true,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Partial Migration Failure",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("ALTER TABLE").WillReturnError(errors.New("insufficient privileges"))
			},
			concurrent:  false,
			wantErr:     true,
			expectedErr: errors.New("insufficient privileges"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to open gorm connection: %v", err)
			}
			defer gormDB.Close()

			tc.setupMock(mock)

			if tc.concurrent {
				var wg sync.WaitGroup
				errChan := make(chan error, 3)

				for i := 0; i < 3; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						err := AutoMigrate(gormDB)
						errChan <- err
					}()
				}

				wg.Wait()
				close(errChan)

				for err := range errChan {
					if (err != nil) != tc.wantErr {
						t.Errorf("Concurrent AutoMigrate() error = %v, wantErr %v", err, tc.wantErr)
					}
				}
			} else {
				err := AutoMigrate(gormDB)
				if (err != nil) != tc.wantErr {
					t.Errorf("AutoMigrate() error = %v, wantErr %v", err, tc.wantErr)
				}
				if tc.wantErr && err != nil && err.Error() != tc.expectedErr.Error() {
					t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}
