package store

import (
	"errors"
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreGetArticles(t *testing.T) {

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
		tagName     string
		username    string
		favoritedBy *model.User
		limit       int64
		offset      int64
		mockSetup   func(sqlmock.Sqlmock)
		wantErr     bool
		expectedLen int
	}{
		{
			name: "Scenario 1: Get Articles Without Filters",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, "Test Article", "Description", "Body", 1, 0)
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("^SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			expectedLen: 1,
		},
		{
			name:     "Scenario 2: Get Articles By Username",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, "Test Article", "Description", "Body", 1, 0)
				mock.ExpectQuery("^SELECT (.+) FROM `articles` join users").
					WithArgs("testuser").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("^SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			expectedLen: 1,
		},
		{
			name:    "Scenario 3: Get Articles By Tag",
			tagName: "testtag",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, "Test Article", "Description", "Body", 1, 0)
				mock.ExpectQuery("^SELECT (.+) FROM `articles` join article_tags").
					WithArgs("testtag").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("^SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			expectedLen: 1,
		},
		{
			name: "Scenario 4: Get Favorited Articles",
			favoritedBy: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {

				favRows := sqlmock.NewRows([]string{"article_id"}).
					AddRow(1)
				mock.ExpectQuery("^SELECT article_id FROM `favorite_articles`").
					WithArgs(1).
					WillReturnRows(favRows)

				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, "Test Article", "Description", "Body", 1, 0)
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("^SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			expectedLen: 1,
		},
		{
			name: "Scenario 6: Handle Empty Result Set",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id", "favorites_count"})
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WillReturnRows(rows)
			},
			expectedLen: 0,
		},
		{
			name: "Scenario 7: Handle Database Error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup(mock)

			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, articles, tt.expectedLen)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
