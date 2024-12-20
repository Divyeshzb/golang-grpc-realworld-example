package store

import (
	"errors"
	"sync"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f


 */
func TestArticleStoreDeleteFavorite(t *testing.T) {
	type testCase struct {
		name          string
		article       *model.Article
		user          *model.User
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
		expectedCount int32
	}

	baseArticle := &model.Article{
		Title:          "Test Article",
		Description:    "Test Description",
		Body:           "Test Body",
		FavoritesCount: 1,
	}
	baseArticle.ID = 1

	baseUser := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}
	baseUser.ID = 1

	tests := []testCase{
		{
			name:    "Successful Deletion",
			article: baseArticle,
			user:    baseUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(baseUser.ID, baseArticle.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles` SET `favorites_count` = favorites_count - \\? WHERE `id` = \\?").
					WithArgs(1, baseArticle.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
			expectedCount: 0,
		},
		{
			name:    "Failed Association Deletion",
			article: baseArticle,
			user:    baseUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(baseUser.ID, baseArticle.ID).
					WillReturnError(errors.New("association deletion failed"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("association deletion failed"),
			expectedCount: 1,
		},
		{
			name:    "Failed FavoritesCount Update",
			article: baseArticle,
			user:    baseUser,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorite_articles`").
					WithArgs(baseUser.ID, baseArticle.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles` SET `favorites_count` = favorites_count - \\? WHERE `id` = \\?").
					WithArgs(1, baseArticle.ID).
					WillReturnError(errors.New("update failed"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("update failed"),
			expectedCount: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			gormDB.LogMode(true)
			tc.setupMock(mock)

			store := &ArticleStore{db: gormDB}

			startCount := tc.article.FavoritesCount
			err = store.DeleteFavorite(tc.article, tc.user)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Equal(t, startCount, tc.article.FavoritesCount)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, tc.article.FavoritesCount)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}

	t.Run("Concurrent Deletions", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		defer db.Close()

		gormDB, err := gorm.Open("mysql", db)
		if err != nil {
			t.Fatalf("Failed to create GORM DB: %v", err)
		}
		defer gormDB.Close()

		article := &model.Article{
			FavoritesCount: 3,
		}
		article.ID = 1

		users := []*model.User{
			{Username: "user1", Email: "user1@example.com", Password: "pass1"},
			{Username: "user2", Email: "user2@example.com", Password: "pass2"},
			{Username: "user3", Email: "user3@example.com", Password: "pass3"},
		}
		for i := range users {
			users[i].ID = uint(i + 1)
		}

		for _, user := range users {
			mock.ExpectBegin()
			mock.ExpectExec("DELETE FROM `favorite_articles`").
				WithArgs(user.ID, article.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("UPDATE `articles` SET `favorites_count` = favorites_count - \\? WHERE `id` = \\?").
				WithArgs(1, article.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
		}

		store := &ArticleStore{db: gormDB}

		var wg sync.WaitGroup
		var mu sync.Mutex
		for _, user := range users {
			wg.Add(1)
			go func(u *model.User) {
				defer wg.Done()
				mu.Lock()
				err := store.DeleteFavorite(article, u)
				mu.Unlock()
				assert.NoError(t, err)
			}(user)
		}
		wg.Wait()

		assert.Equal(t, int32(0), article.FavoritesCount)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %s", err)
		}
	})
}

