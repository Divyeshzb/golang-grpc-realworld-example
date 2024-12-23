package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		mockDB  func(mock sqlmock.Sqlmock, user *model.User)
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful Update",
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "updated_user",
				Email:    "updated@example.com",
				Password: "newpassword123",
				Bio:      "Updated bio",
				Image:    "updated.jpg",
			},
			mockDB: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users` SET").
					WithArgs(
						user.Username,
						user.Email,
						user.Password,
						user.Bio,
						user.Image,
						user.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Update Non-Existent User",
			user: &model.User{
				Model:    gorm.Model{ID: 999},
				Username: "nonexistent",
				Email:    "nonexistent@example.com",
			},
			mockDB: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users` SET").
					WithArgs(
						user.Username,
						user.Email,
						user.ID,
					).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "record not found",
		},
		{
			name: "Update with Duplicate Username",
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "duplicate_user",
				Email:    "unique@example.com",
			},
			mockDB: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users` SET").
					WithArgs(
						user.Username,
						user.Email,
						user.ID,
					).
					WillReturnError(errors.New("Error 1062: Duplicate entry"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "Duplicate entry",
		},
		{
			name:    "Update with Nil User",
			user:    nil,
			mockDB:  func(mock sqlmock.Sqlmock, user *model.User) {},
			wantErr: true,
			errMsg:  "invalid user object",
		},
		{
			name: "Database Connection Error",
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "test_user",
				Email:    "test@example.com",
			},
			mockDB: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users` SET").
					WithArgs(
						user.Username,
						user.Email,
						user.ID,
					).
					WillReturnError(errors.New("database connection error"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "database connection error",
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
				t.Fatalf("Failed to open gorm DB: %v", err)
			}
			defer gormDB.Close()

			gormDB.LogMode(true)
			store := &UserStore{db: gormDB}

			if tt.user != nil {
				tt.mockDB(mock, tt.user)
			}

			err = store.Update(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
