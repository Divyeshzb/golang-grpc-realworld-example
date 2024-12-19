package store

import (
	"errors"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreUpdate(t *testing.T) {

	type testCase struct {
		name          string
		user          *model.User
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successful Update",
			user: &model.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Database Error",
			user: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
				Username: "testuser",
				Email:    "test@example.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("database error"),
		},
		{
			name: "Nil User",
			user: nil,
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectedError: errors.New("invalid user model"),
		},
		{
			name: "Unique Constraint Violation",
			user: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
				Username: "existinguser",
				Email:    "existing@example.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnError(errors.New("Error 1062: Duplicate entry"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("Error 1062: Duplicate entry"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			if tc.setupMock != nil {
				tc.setupMock(mock)
			}

			store := &UserStore{
				db: gormDB,
			}

			err = store.Update(tc.user)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test case '%s' completed", tc.name)
		})
	}
}
