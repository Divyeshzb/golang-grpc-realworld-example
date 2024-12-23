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
			name: "Successful deletion",
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				Title:          "Test Article",
				Description:    "Test Description",
				Body:           "Test Body",
				FavoritesCount: 2,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Failed association deletion",
			article: &model.Article{
				Model:          gorm.Model{ID: 2},
				Title:          "Test Article 2",
				Description:    "Test Description 2",
				Body:           "Test Body 2",
				FavoritesCount: 1,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "testuser2",
				Email:    "test2@example.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(2, 2).
					WillReturnError(errors.New("association deletion failed"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "association deletion failed",
		},
		{
			name: "Failed favorites count update",
			article: &model.Article{
				Model:          gorm.Model{ID: 3},
				Title:          "Test Article 3",
				Description:    "Test Description 3",
				Body:           "Test Body 3",
				FavoritesCount: 1,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 3},
				Username: "testuser3",
				Email:    "test3@example.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(3, 3).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles` SET").
					WithArgs(1, 3).
					WillReturnError(errors.New("update failed"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.setupMock(mock)

			initialCount := tt.article.FavoritesCount

			err := store.DeleteFavorite(tt.article, tt.user)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Equal(t, initialCount, tt.article.FavoritesCount)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, initialCount-1, tt.article.FavoritesCount)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
