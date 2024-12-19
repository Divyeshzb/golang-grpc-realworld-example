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

func TestUserStoreUnfollow(t *testing.T) {

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

	store := &UserStore{db: gormDB}

	tests := []struct {
		name        string
		setupMock   func()
		userA       *model.User
		userB       *model.User
		expectedErr error
	}{
		{
			name: "Successful Unfollow",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
			},
			expectedErr: nil,
		},
		{
			name: "Unfollow Non-Followed User",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 3).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 3},
				Username: "userC",
				Email:    "userC@test.com",
			},
			expectedErr: nil,
		},
		{
			name: "Invalid User References",
			setupMock: func() {

			},
			userA:       nil,
			userB:       nil,
			expectedErr: errors.New("invalid user reference"),
		},
		{
			name: "Database Connection Error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("database connection error"))
				mock.ExpectRollback()
			},
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
			},
			expectedErr: errors.New("database connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.setupMock != nil {
				tt.setupMock()
			}

			err := store.Unfollow(tt.userA, tt.userB)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}

	t.Run("Concurrent Unfollow Operations", func(t *testing.T) {
		userA := &model.User{
			Model:    gorm.Model{ID: 1},
			Username: "userA",
			Email:    "userA@test.com",
		}
		userB := &model.User{
			Model:    gorm.Model{ID: 2},
			Username: "userB",
			Email:    "userB@test.com",
		}

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM `follows`").
			WithArgs(1, 2).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		done := make(chan bool)
		go func() {
			err := store.Unfollow(userA, userB)
			assert.NoError(t, err)
			done <- true
		}()

		select {
		case <-done:

		case <-time.After(2 * time.Second):
			t.Error("Concurrent operation timed out")
		}
	})
}
