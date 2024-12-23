package store

import (
	"database/sql"
	"testing"
	"time"
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
		expected    []model.Article
		expectError bool
	}{
		{
			name:   "Scenario 1: Get Articles Without Filters",
			limit:  10,
			offset: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article", "Description", "Body", 1)
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("^SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			expected: []model.Article{
				{
					Model: gorm.Model{ID: 1},
					Title: "Test Article",
					Author: model.User{
						Model:    gorm.Model{ID: 1},
						Username: "testuser",
						Email:    "test@example.com",
					},
				},
			},
		},
		{
			name:     "Scenario 2: Get Articles By Username",
			username: "testuser",
			limit:    10,
			offset:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "User Article", 1)
				mock.ExpectQuery("^SELECT (.+) FROM `articles` join users").
					WithArgs("testuser").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username"}).
					AddRow(1, "testuser")
				mock.ExpectQuery("^SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			expected: []model.Article{
				{
					Model: gorm.Model{ID: 1},
					Title: "User Article",
					Author: model.User{
						Model:    gorm.Model{ID: 1},
						Username: "testuser",
					},
				},
			},
		},
		{
			name:    "Scenario 3: Get Articles By Tag",
			tagName: "programming",
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Programming Article", 1)
				mock.ExpectQuery("^SELECT (.+) FROM `articles` join article_tags").
					WithArgs("programming").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username"}).
					AddRow(1, "testuser")
				mock.ExpectQuery("^SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			expected: []model.Article{
				{
					Model: gorm.Model{ID: 1},
					Title: "Programming Article",
					Author: model.User{
						Model:    gorm.Model{ID: 1},
						Username: "testuser",
					},
				},
			},
		},
		{
			name: "Scenario 4: Get Favorited Articles",
			favoritedBy: &model.User{
				Model: gorm.Model{ID: 1},
			},
			limit:  10,
			offset: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				favRows := sqlmock.NewRows([]string{"article_id"}).
					AddRow(1)
				mock.ExpectQuery("^SELECT article_id FROM `favorite_articles`").
					WithArgs(1).
					WillReturnRows(favRows)

				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "user_id"}).
					AddRow(1, time.Now(), time.Now(), nil, "Favorited Article", 1)
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username"}).
					AddRow(1, "testuser")
				mock.ExpectQuery("^SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			expected: []model.Article{
				{
					Model: gorm.Model{ID: 1},
					Title: "Favorited Article",
					Author: model.User{
						Model:    gorm.Model{ID: 1},
						Username: "testuser",
					},
				},
			},
		},
		{
			name:   "Scenario 8: Database Error",
			limit:  10,
			offset: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WillReturnError(sql.ErrConnDone)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup(mock)

			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expected), len(articles))
				if len(articles) > 0 {
					assert.Equal(t, tt.expected[0].ID, articles[0].ID)
					assert.Equal(t, tt.expected[0].Title, articles[0].Title)
					assert.Equal(t, tt.expected[0].Author.Username, articles[0].Author.Username)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
