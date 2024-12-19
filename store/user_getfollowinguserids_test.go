package store

import (
	"database/sql"
	"testing"
	"time"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreGetFollowingUserIDs(t *testing.T) {

	type testCase struct {
		name          string
		userID        uint
		mockSetup     func(sqlmock.Sqlmock)
		expectedIDs   []uint
		expectedError error
	}

	tests := []testCase{
		{
			name:   "Successful retrieval of following user IDs",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(2).
					AddRow(3).
					AddRow(4)
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{2, 3, 4},
			expectedError: nil,
		},
		{
			name:   "User with no followings",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{},
			expectedError: nil,
		},
		{
			name:   "Database error",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedIDs:   []uint{},
			expectedError: sql.ErrConnDone,
		},
		{
			name:   "Large number of followings",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				for i := uint(2); i <= 1001; i++ {
					rows.AddRow(i)
				}
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs: func() []uint {
				ids := make([]uint, 1000)
				for i := uint(0); i < 1000; i++ {
					ids[i] = i + 2
				}
				return ids
			}(),
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

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

			tc.mockSetup(mock)

			store := &UserStore{db: gormDB}

			user := &model.User{
				Model: gorm.Model{
					ID:        tc.userID,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
				Bio:      "test bio",
				Image:    "test.jpg",
			}

			ids, err := store.GetFollowingUserIDs(user)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedIDs, ids)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed successfully", tc.name)
		})
	}
}
