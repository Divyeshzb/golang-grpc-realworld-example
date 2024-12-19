package store

import (
	"database/sql"
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreCreate(t *testing.T) {

	tests := []struct {
		name    string
		user    *model.User
		wantErr bool
		setup   func(mock sqlmock.Sqlmock)
	}{
		{
			name: "Successful User Creation",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			wantErr: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "Duplicate Username",
			user: &model.User{
				Username: "existinguser",
				Email:    "new@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			wantErr: true,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
		},
		{
			name: "Duplicate Email",
			user: &model.User{
				Username: "newuser",
				Email:    "existing@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			wantErr: true,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
		},
		{
			name: "Missing Required Fields",
			user: &model.User{
				Username: "",
				Email:    "",
				Password: "",
				Bio:      "",
				Image:    "",
			},
			wantErr: true,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
		},
		{
			name: "Maximum Field Lengths",
			user: &model.User{
				Username: string(make([]byte, 255)),
				Email:    "very.long.email@example.com",
				Password: string(make([]byte, 255)),
				Bio:      string(make([]byte, 1000)),
				Image:    string(make([]byte, 255)),
			},
			wantErr: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "Special Characters",
			user: &model.User{
				Username: "user@#$%^",
				Email:    "special.chars+test@example.com",
				Password: "pass!@#$%^&*()",
				Bio:      "Bio with émojis 🎉",
				Image:    "image-with-späces.jpg",
			},
			wantErr: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
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
				t.Fatalf("Failed to open GORM DB: %v", err)
			}
			defer gormDB.Close()

			tt.setup(mock)

			store := &UserStore{
				db: gormDB,
			}

			err = store.Create(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				t.Logf("User created successfully")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
