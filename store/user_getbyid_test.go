package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreGetByID(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name: "Successfully Retrieve User by Valid ID",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("SELECT (.+) FROM `users` WHERE (.+)").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectedError: nil,
		},
		{
			name: "Non-existent User ID",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `users` WHERE (.+)").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `users` WHERE (.+)").
					WithArgs(1).
					WillReturnError(errors.New("database connection error"))
			},
			expectedUser:  nil,
			expectedError: errors.New("database connection error"),
		},
		{
			name: "Zero ID Input",
			id:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `users` WHERE (.+)").
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Partial Data Return",
			id:   2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(2, "partialuser", nil)
				mock.ExpectQuery("SELECT (.+) FROM `users` WHERE (.+)").
					WithArgs(2).
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "partialuser",
			},
			expectedError: nil,
		},
		{
			name: "Maximum ID Value",
			id:   ^uint(0),
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `users` WHERE (.+)").
					WithArgs(^uint(0)).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database connection: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB instance: %v", err)
			}
			defer gormDB.Close()

			tt.mockSetup(mock)

			userStore := &UserStore{db: gormDB}

			user, err := userStore.GetByID(tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.Model.ID, user.Model.ID)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed successfully", tt.name)
		})
	}
}
