package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreDeleteFavorite(t *testing.T) {

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

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name        string
		article     *model.Article
		user        *model.User
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successful Deletion",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 2,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Failed Association Deletion",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 2,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WillReturnError(errors.New("association deletion failed"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "association deletion failed",
		},
		{
			name: "Failed FavoritesCount Update",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 2,
			},
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(errors.New("update failed"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "update failed",
		},
		{
			name:    "Nil Article",
			article: nil,
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectError: true,
			errorMsg:    "invalid transaction",
		},
		{
			name: "Nil User",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				FavoritesCount: 2,
			},
			user: nil,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WillReturnError(errors.New("invalid user"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "invalid user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.setupMock(mock)

			err := store.DeleteFavorite(tt.article, tt.user)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				if tt.article != nil {
					assert.Equal(t, tt.article.FavoritesCount-1, tt.article.FavoritesCount)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
