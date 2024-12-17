package store

import (
	"database/sql"
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreGetArticles(t *testing.T) {

	setupTestDB := func(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *ArticleStore) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}

		gormDB, err := gorm.Open("mysql", db)
		if err != nil {
			t.Fatalf("Failed to open GORM DB: %v", err)
		}

		return db, mock, &ArticleStore{db: gormDB}
	}

	tests := []struct {
		name        string
		tagName     string
		username    string
		favoritedBy *model.User
		limit       int64
		offset      int64
		setupMock   func(sqlmock.Sqlmock)
		wantErr     bool
		wantLen     int
	}{
		{
			name: "Scenario 1: Get Articles Without Filters",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, "Test Article", "Description", "Body", 1, 0).
					AddRow(2, "Test Article 2", "Description 2", "Body 2", 1, 0)

				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WillReturnRows(rows)
			},
			limit:   10,
			wantLen: 2,
			wantErr: false,
		},
		{
			name:     "Scenario 2: Get Articles By Username",
			username: "testuser",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, "Test Article", "Description", "Body", 1, 0)

				mock.ExpectQuery("^SELECT (.+) FROM `articles` join users").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			limit:   10,
			wantLen: 1,
			wantErr: false,
		},
		{
			name:    "Scenario 3: Get Articles By Tag",
			tagName: "testtag",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, "Test Article", "Description", "Body", 1, 0)

				mock.ExpectQuery("^SELECT (.+) FROM `articles` join article_tags").
					WithArgs("testtag").
					WillReturnRows(rows)
			},
			limit:   10,
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "Scenario 4: Get Favorited Articles",
			favoritedBy: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				favRows := sqlmock.NewRows([]string{"article_id"}).
					AddRow(1)

				mock.ExpectQuery("^SELECT article_id FROM `favorite_articles`").
					WithArgs(1).
					WillReturnRows(favRows)

				articleRows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id", "favorites_count"}).
					AddRow(1, "Test Article", "Description", "Body", 1, 0)

				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WillReturnRows(articleRows)
			},
			limit:   10,
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "Scenario 6: Empty Result Set",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id", "favorites_count"})
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WillReturnRows(rows)
			},
			limit:   10,
			wantLen: 0,
			wantErr: false,
		},
		{
			name: "Scenario 7: Database Error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WillReturnError(sql.ErrConnDone)
			},
			limit:   10,
			wantErr: true,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, store := setupTestDB(t)
			defer db.Close()

			tt.setupMock(mock)

			articles, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, articles, tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
