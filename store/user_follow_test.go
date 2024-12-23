package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestFollow(t *testing.T) {

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
		follower    *model.User
		following   *model.User
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successful Follow",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			follower: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "user1",
				Email:    "user1@test.com",
			},
			following: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "user2",
				Email:    "user2@test.com",
			},
			expectError: false,
		},
		{
			name: "Follow Already Followed User",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("duplicate entry"))
				mock.ExpectRollback()
			},
			follower: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "user1",
				Email:    "user1@test.com",
			},
			following: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "user2",
				Email:    "user2@test.com",
			},
			expectError: true,
			errorMsg:    "duplicate entry",
		},
		{
			name: "Self-Follow Attempt",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			follower: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "user1",
				Email:    "user1@test.com",
			},
			following: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "user1",
				Email:    "user1@test.com",
			},
			expectError: false,
		},
		{
			name: "Database Connection Error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("connection refused"))
				mock.ExpectRollback()
			},
			follower: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "user1",
				Email:    "user1@test.com",
			},
			following: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "user2",
				Email:    "user2@test.com",
			},
			expectError: true,
			errorMsg:    "connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.setupMock()

			err := store.Follow(tt.follower, tt.following)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
