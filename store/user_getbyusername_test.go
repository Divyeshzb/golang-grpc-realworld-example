package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreGetByUsername(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:     "Successfully retrieve user by valid username",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"username", "email"}).
					AddRow("testuser", "test@example.com")
				mock.ExpectQuery("^SELECT (.+) FROM \"users\"").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectedError: nil,
		},
		{
			name:     "Handle non-existent username",
			username: "nonexistent",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\"").
					WithArgs("nonexistent").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Handle database connection error",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\"").
					WithArgs("testuser").
					WillReturnError(errors.New("database connection error"))
			},
			expectedUser:  nil,
			expectedError: errors.New("database connection error"),
		},
		{
			name:     "Handle empty username parameter",
			username: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\"").
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Handle special characters in username",
			username: "test@user#123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"username", "email"}).
					AddRow("test@user#123", "special@example.com")
				mock.ExpectQuery("^SELECT (.+) FROM \"users\"").
					WithArgs("test@user#123").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Username: "test@user#123",
				Email:    "special@example.com",
			},
			expectedError: nil,
		},
		{
			name:     "Handle case sensitivity in username",
			username: "TestUser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"username", "email"}).
					AddRow("TestUser", "case@example.com")
				mock.ExpectQuery("^SELECT (.+) FROM \"users\"").
					WithArgs("TestUser").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Username: "TestUser",
				Email:    "case@example.com",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database connection: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB instance: %v", err)
			}
			defer gormDB.Close()

			tt.mockSetup(mock)

			userStore := &UserStore{
				db: gormDB,
			}

			user, err := userStore.GetByUsername(tt.username)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
			}

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
			}

			t.Logf("Test case '%s' completed successfully", tt.name)
		})
	}
}
