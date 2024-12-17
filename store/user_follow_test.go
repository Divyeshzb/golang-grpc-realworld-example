package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreFollow(t *testing.T) {

	type testCase struct {
		name    string
		userA   *model.User
		userB   *model.User
		setupDB func(mock sqlmock.Sqlmock)
		wantErr bool
	}

	setupMockDB := func() (*gorm.DB, sqlmock.Sqlmock) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}

		gormDB, err := gorm.Open("mysql", db)
		if err != nil {
			t.Fatalf("Failed to open GORM DB: %v", err)
		}

		return gormDB, mock
	}

	tests := []testCase{
		{
			name: "Successful Follow Operation",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
				Password: "password",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
				Password: "password",
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Follow Already Followed User",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
				Password: "password",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
				Password: "password",
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("duplicate entry"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Self-Follow Attempt",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
				Password: "password",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
				Password: "password",
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 1).
					WillReturnError(errors.New("self-follow not allowed"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:  "Follow with Nil User",
			userA: nil,
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
				Password: "password",
			},
			setupDB: func(mock sqlmock.Sqlmock) {

			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			gormDB, mock := setupMockDB()
			defer gormDB.Close()

			store := &UserStore{db: gormDB}

			if tc.setupDB != nil {
				tc.setupDB(mock)
			}

			err := store.Follow(tc.userA, tc.userB)

			if tc.wantErr {
				assert.Error(t, err)
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				t.Log("Follow operation successful")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
