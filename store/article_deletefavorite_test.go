package store

import (
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"errors"
	"sync"
	"testing"
)

func TestArticleStoreDeleteFavorite(t *testing.T) {

	type testCase struct {
		name          string
		article       *model.Article
		user          *model.User
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
	}

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

	store := &ArticleStore{db: gormDB}

	testUser := &model.User{
		Model:    gorm.Model{ID: 1},
		Username: "testuser",
		Email:    "test@example.com",
	}

	testArticle := &model.Article{
		Model:          gorm.Model{ID: 1},
		Title:          "Test Article",
		FavoritesCount: 1,
	}

	tests := []testCase{
		{
			name:    "Successful Deletion",
			article: testArticle,
			user:    testUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(testArticle.ID, testUser.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name:    "Failed Association Deletion",
			article: testArticle,
			user:    testUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WillReturnError(errors.New("association deletion failed"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("association deletion failed"),
		},
		{
			name:    "Failed FavoritesCount Update",
			article: testArticle,
			user:    testUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(errors.New("update failed"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("update failed"),
		},
		{
			name:    "Nil Article",
			article: nil,
			user:    testUser,
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectedError: errors.New("invalid article"),
		},
		{
			name:    "Nil User",
			article: testArticle,
			user:    nil,
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectedError: errors.New("invalid user"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock(mock)

			err := store.DeleteFavorite(tc.article, tc.user)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				if tc.article != nil {
					assert.Equal(t, tc.article.FavoritesCount, int32(0))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}

	t.Run("Concurrent Operations", func(t *testing.T) {
		const numGoroutines = 5
		var wg sync.WaitGroup
		articleWithMultipleFavorites := &model.Article{
			Model:          gorm.Model{ID: 2},
			Title:          "Concurrent Test Article",
			FavoritesCount: int32(numGoroutines),
		}

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM `favorite_articles`").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("UPDATE `articles`").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(userID uint) {
				defer wg.Done()
				concurrentUser := &model.User{Model: gorm.Model{ID: userID}}
				_ = store.DeleteFavorite(articleWithMultipleFavorites, concurrentUser)
			}(uint(i + 1))
		}

		wg.Wait()
	})
}
