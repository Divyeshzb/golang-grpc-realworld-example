package store

import (
	"database/sql"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}
	defer gormDB.Close()

	userStore := &UserStore{db: gormDB}

	tests := []struct {
		name        string
		user        *model.User
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successful user creation",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Duplicate username",
			user: &model.User{
				Username: "existinguser",
				Email:    "new@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "sql: no rows in result set",
		},
		{
			name: "Empty required fields",
			user: &model.User{
				Username: "",
				Email:    "",
				Password: "",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "sql: no rows in result set",
		},
		{
			name: "Database connection error",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "sql: connection is already closed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			err := userStore.Create(tt.user)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Equal(t, tt.errorMsg, err.Error())
				}
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				t.Log("User created successfully")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
