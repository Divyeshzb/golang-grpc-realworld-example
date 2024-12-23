package store

import (
	"database/sql"
	"regexp"
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestGetByUsername(t *testing.T) {
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
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}).
					AddRow(1, "2023-01-01 00:00:00", "2023-01-01 00:00:00", nil, "testuser", "test@example.com", "hashedpassword", "test bio", "image.jpg")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashedpassword",
				Bio:      "test bio",
				Image:    "image.jpg",
			},
			expectedError: nil,
		},
		{
			name:     "Non-existent username",
			username: "nonexistentuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("nonexistentuser").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Empty username",
			username: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Database connection error",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("testuser").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
		{
			name:     "Special characters in username",
			username: "test@user_123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}).
					AddRow(1, "2023-01-01 00:00:00", "2023-01-01 00:00:00", nil, "test@user_123", "special@example.com", "hashedpassword", "test bio", "image.jpg")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) ORDER BY `users`.`id` ASC LIMIT 1")).
					WithArgs("test@user_123").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "test@user_123",
				Email:    "special@example.com",
				Password: "hashedpassword",
				Bio:      "test bio",
				Image:    "image.jpg",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			tt.mockSetup(mock)

			userStore := &UserStore{db: gormDB}
			user, err := userStore.GetByUsername(tt.username)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Bio, user.Bio)
				assert.Equal(t, tt.expectedUser.Image, user.Image)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
