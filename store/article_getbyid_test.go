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

func TestArticleStoreGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name          string
		id            uint
		mockSetup     func(sqlmock.Sqlmock)
		expectedError error
		expectedData  *model.Article
	}{
		{
			name: "Successful retrieval",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "description", "body"}).
					AddRow(1, time.Now(), time.Now(), nil, "Test Article", "Test Description", "Test Body")
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(1).
					WillReturnRows(rows)

				tagRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
					AddRow(1, time.Now(), time.Now(), nil, "tag1").
					AddRow(2, time.Now(), time.Now(), nil, "tag2")
				mock.ExpectQuery(`SELECT \* FROM "tags"`).
					WillReturnRows(tagRows)

				authorRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email"}).
					AddRow(1, time.Now(), time.Now(), nil, "testuser", "test@example.com")
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WillReturnRows(authorRows)
			},
			expectedError: nil,
			expectedData: &model.Article{
				Model: gorm.Model{
					ID: 1,
				},
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				Tags: []model.Tag{
					{Model: gorm.Model{ID: 1}, Name: "tag1"},
					{Model: gorm.Model{ID: 2}, Name: "tag2"},
				},
				Author: model.User{
					Model:    gorm.Model{ID: 1},
					Username: "testuser",
					Email:    "test@example.com",
				},
			},
		},
		{
			name: "Non-existent article",
			id:   99999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(99999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
			expectedData:  nil,
		},
		{
			name: "Database error",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: sql.ErrConnDone,
			expectedData:  nil,
		},
		{
			name: "Zero ID",
			id:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			article, err := store.GetByID(tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, article)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, article)
				assert.Equal(t, tt.expectedData.ID, article.ID)
				assert.Equal(t, tt.expectedData.Title, article.Title)
				assert.Equal(t, tt.expectedData.Description, article.Description)
				assert.Equal(t, tt.expectedData.Body, article.Body)
				assert.Len(t, article.Tags, len(tt.expectedData.Tags))
				assert.Equal(t, tt.expectedData.Author.ID, article.Author.ID)
				assert.Equal(t, tt.expectedData.Author.Username, article.Author.Username)
				assert.Equal(t, tt.expectedData.Author.Email, article.Author.Email)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}
