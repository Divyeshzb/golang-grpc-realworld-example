package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreAddFavorite(t *testing.T) {
	type testCase struct {
		name          string
		article       *model.Article
		user          *model.User
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successfully Add Favorite",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Handle Database Transaction Failure During User Association",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WithArgs(1, 1).
					WillReturnError(errors.New("association error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("association error"),
		},
		{
			name: "Handle Database Transaction Failure During Count Update",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(1, 1).
					WillReturnError(errors.New("update error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("update error"),
		},
		{
			name:    "Add Favorite with Nil Article",
			article: nil,
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock:     func(mock sqlmock.Sqlmock) {},
			expectedError: errors.New("invalid article"),
		},
		{
			name: "Add Favorite with Nil User",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 0,
			},
			user:          nil,
			setupMock:     func(mock sqlmock.Sqlmock) {},
			expectedError: errors.New("invalid user"),
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

			gormDB.LogMode(false)

			store := &ArticleStore{
				db: gormDB,
			}

			if tc.article == nil && tc.expectedError.Error() == "invalid article" {
				err = store.AddFavorite(tc.article, tc.user)
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				return
			}

			if tc.user == nil && tc.expectedError.Error() == "invalid user" {
				err = store.AddFavorite(tc.article, tc.user)
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				return
			}

			if tc.setupMock != nil {
				tc.setupMock(mock)
			}

			err = store.AddFavorite(tc.article, tc.user)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				if tc.article != nil {
					assert.Equal(t, int32(1), tc.article.FavoritesCount)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
