package store

import (
	"database/sql"
	"errors"
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestGetFollowingUserIDs(t *testing.T) {
	tests := []struct {
		name          string
		user          *model.User
		mockSetup     func(sqlmock.Sqlmock)
		expectedIDs   []uint
		expectedError error
	}{
		{
			name: "Successfully Retrieve Following User IDs",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(2).
					AddRow(3).
					AddRow(4)
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE .*from_user_id = \\?.*$").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{2, 3, 4},
			expectedError: nil,
		},
		{
			name: "User with No Followings",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE .*from_user_id = \\?.*$").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{},
			expectedError: nil,
		},
		{
			name: "Database Connection Error",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE .*from_user_id = \\?.*$").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			expectedIDs:   []uint{},
			expectedError: errors.New("database error"),
		},
		{
			name: "Invalid User ID",
			user: &model.User{Model: gorm.Model{ID: 0}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE .*from_user_id = \\?.*$").
					WithArgs(0).
					WillReturnError(sql.ErrNoRows)
			},
			expectedIDs:   []uint{},
			expectedError: sql.ErrNoRows,
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

			store := &UserStore{db: gormDB}
			ids, err := store.GetFollowingUserIDs(tt.user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedIDs, ids)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
