package store

import (
	"database/sql"
	"math"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	createdAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

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
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE \\(`users`.`id` = \\?\\) AND `users`.`deleted_at` IS NULL").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, createdAt, updatedAt, nil, "testuser", "test@example.com", "hashedpass", "test bio", "image.jpg"))
			},
			expectedUser: &model.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashedpass",
				Bio:      "test bio",
				Image:    "image.jpg",
			},
			expectedError: nil,
		},
		{
			name: "Attempt to Retrieve Non-existent User ID",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE \\(`users`.`id` = \\?\\) AND `users`.`deleted_at` IS NULL").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Handle Zero ID Value",
			id:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE \\(`users`.`id` = \\?\\) AND `users`.`deleted_at` IS NULL").
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE \\(`users`.`id` = \\?\\) AND `users`.`deleted_at` IS NULL").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
		{
			name: "Handle Maximum uint ID Value",
			id:   math.MaxUint32,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE \\(`users`.`id` = \\?\\) AND `users`.`deleted_at` IS NULL").
					WithArgs(math.MaxUint32).
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
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			gormDB.LogMode(true)
			defer gormDB.Close()

			tt.mockSetup(mock)

			store := &UserStore{db: gormDB}
			user, err := store.GetByID(tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Bio, user.Bio)
				assert.Equal(t, tt.expectedUser.Image, user.Image)
				assert.Equal(t, tt.expectedUser.CreatedAt, user.CreatedAt)
				assert.Equal(t, tt.expectedUser.UpdatedAt, user.UpdatedAt)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
