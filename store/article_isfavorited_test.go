package store

import (
	"errors"
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreIsFavorited(t *testing.T) {

	type testCase struct {
		name        string
		article     *model.Article
		user        *model.User
		mockSetup   func(sqlmock.Sqlmock)
		expected    bool
		expectedErr error
	}

	tests := []testCase{
		{
			name: "Scenario 1: Valid Article and User with Existing Favorite Relationship",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorite_articles"`).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:    true,
			expectedErr: nil,
		},
		{
			name: "Scenario 2: Valid Article and User with No Favorite Relationship",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorite_articles"`).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Scenario 3: Nil Article Parameter",
			article:     nil,
			user:        &model.User{Model: gorm.Model{ID: 1}},
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Scenario 4: Nil User Parameter",
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			user:        nil,
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name: "Scenario 5: Database Error Condition",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorite_articles"`).
					WithArgs(1, 1).
					WillReturnError(errors.New("database error"))
			},
			expected:    false,
			expectedErr: errors.New("database error"),
		},
		{
			name: "Scenario 6: Zero ID Values",
			article: &model.Article{
				Model: gorm.Model{ID: 0},
			},
			user: &model.User{
				Model: gorm.Model{ID: 0},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorite_articles"`).
					WithArgs(0, 0).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:    false,
			expectedErr: nil,
		},
		{
			name: "Scenario 7: Multiple Favorite Relationships",
			article: &model.Article{
				Model: gorm.Model{ID: 1},
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorite_articles"`).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
			},
			expected:    true,
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.mockSetup(mock)

			store := &ArticleStore{db: gormDB}

			result, err := store.IsFavorited(tc.article, tc.user)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, result)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
		})
	}
}
