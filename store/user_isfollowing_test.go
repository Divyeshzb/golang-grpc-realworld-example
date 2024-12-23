package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreIsFollowing(t *testing.T) {

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

	userStore := &UserStore{db: gormDB}

	tests := []struct {
		name        string
		userA       *model.User
		userB       *model.User
		mockSetup   func(sqlmock.Sqlmock)
		expected    bool
		expectedErr error
	}{
		{
			name: "Valid Following Relationship",
			userA: &model.User{Model: gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
				Password: "password",
				Bio:      "bio",
				Image:    "image"},
			userB: &model.User{Model: gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
				Password: "password",
				Bio:      "bio",
				Image:    "image"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:    true,
			expectedErr: nil,
		},
		{
			name:  "Valid Non-Following Relationship",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil User A",
			userA:       nil,
			userB:       &model.User{Model: gorm.Model{ID: 2}},
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil User B",
			userA:       &model.User{Model: gorm.Model{ID: 1}},
			userB:       nil,
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Both Users Nil",
			userA:       nil,
			userB:       nil,
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:  "Database Error",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))
			},
			expected:    false,
			expectedErr: errors.New("database error"),
		},
		{
			name:  "Same User Check",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:  "Multiple Following Relationships",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:    true,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup(mock)

			result, err := userStore.IsFollowing(tt.userA, tt.userB)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test case '%s' completed successfully", tt.name)
		})
	}
}
