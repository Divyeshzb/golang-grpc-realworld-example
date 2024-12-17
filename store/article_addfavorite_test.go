package store

import (
	"errors"
	"sync"
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreAddFavorite(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		article     *model.Article
		user        *model.User
		expectError bool
	}{
		{
			name: "Successful Addition of Favorite",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			article: &model.Article{
				Model:          gorm.Model{ID: 1},
				Title:          "Test Article",
				FavoritesCount: 0,
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectError: false,
		},
		{
			name: "Database Error During User Association",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectError: true,
		},
		{
			name: "Database Error During Favorites Count Update",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `favorite_articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			article: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectError: true,
		},
		{
			name: "Null Parameter Handling",
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			article:     nil,
			user:        nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock, err := setupTestDB(t)
			assert.NoError(t, err)
			defer gormDB.Close()

			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			store := &ArticleStore{db: gormDB}

			err = store.AddFavorite(tt.article, tt.user)

			if tt.expectError {
				assert.Error(t, err)
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int32(1), tt.article.FavoritesCount)
				t.Log("Favorite added successfully")
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
func TestArticleStoreAddFavoriteConcurrent(t *testing.T) {
	gormDB, mock, err := setupTestDB(t)
	assert.NoError(t, err)
	defer gormDB.Close()

	article := &model.Article{
		Model:          gorm.Model{ID: 1},
		Title:          "Test Article",
		FavoritesCount: 0,
	}

	numUsers := 5
	var wg sync.WaitGroup
	store := &ArticleStore{db: gormDB}

	mock.ExpectBegin()
	for i := 0; i < numUsers; i++ {
		mock.ExpectExec("INSERT INTO `favorite_articles`").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("UPDATE `articles`").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectBegin()
	}

	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func(userID uint) {
			defer wg.Done()
			user := &model.User{
				Model:    gorm.Model{ID: userID},
				Username: "testuser",
				Email:    "test@example.com",
			}
			err := store.AddFavorite(article, user)
			assert.NoError(t, err)
		}(uint(i + 1))
	}

	wg.Wait()

	assert.Equal(t, int32(numUsers), article.FavoritesCount)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}
