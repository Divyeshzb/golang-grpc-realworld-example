package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreCreate(t *testing.T) {

	tests := []struct {
		name    string
		article *model.Article
		dbErr   error
		wantErr bool
	}{
		{
			name: "Success - Create article with all required fields",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body Content",
				UserID:      1,
			},
			dbErr:   nil,
			wantErr: false,
		},
		{
			name: "Success - Create article with tags",
			article: &model.Article{
				Title:       "Article with Tags",
				Description: "Description with tags",
				Body:        "Body content with tags",
				UserID:      1,
				Tags:        []model.Tag{{Name: "test-tag"}},
			},
			dbErr:   nil,
			wantErr: false,
		},
		{
			name: "Failure - Database error",
			article: &model.Article{
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body Content",
				UserID:      1,
			},
			dbErr:   errors.New("database error"),
			wantErr: true,
		},
		{
			name: "Failure - Missing required fields",
			article: &model.Article{

				Description: "Test Description",
				Body:        "Test Body Content",
				UserID:      1,
			},
			dbErr:   errors.New("title cannot be null"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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

			store := &ArticleStore{
				db: gormDB,
			}

			if !tt.wantErr {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				if len(tt.article.Tags) > 0 {
					mock.ExpectExec("INSERT INTO `article_tags`").
						WillReturnResult(sqlmock.NewResult(1, 1))
				}
				mock.ExpectCommit()
			} else {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnError(tt.dbErr)
				mock.ExpectRollback()
			}

			err = store.Create(tt.article)

			if tt.wantErr {
				assert.Error(t, err)
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.article.ID, "Article ID should be set after creation")
				t.Log("Article created successfully")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
