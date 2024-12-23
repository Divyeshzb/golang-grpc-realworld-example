package store

import (
	"database/sql"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {

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

	now := time.Now()

	tests := []struct {
		name          string
		id            uint
		mockSetup     func(sqlmock.Sqlmock)
		expectedError error
		expectedData  *model.Article
	}{
		{
			name: "Successful retrieval with valid ID",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "author_id",
				}).AddRow(1, now, now, nil, "Test Article", "Test Description", "Test Body", 1)

				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WithArgs(1).
					WillReturnRows(rows)

				tagRows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "tag1").
					AddRow(2, "tag2")
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").
					WillReturnRows(tagRows)

				authorRows := sqlmock.NewRows([]string{
					"id", "username", "email", "password", "bio", "image",
				}).AddRow(1, "testuser", "test@example.com", "hashedpass", "test bio", "test.jpg")
				mock.ExpectQuery("^SELECT (.+) FROM `users`").
					WillReturnRows(authorRows)
			},
			expectedError: nil,
			expectedData: &model.Article{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
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
					Password: "hashedpass",
					Bio:      "test bio",
					Image:    "test.jpg",
				},
			},
		},
		{
			name: "Non-existent article ID",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
			expectedData:  nil,
		},
		{
			name: "Database connection error",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `articles`").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: sql.ErrConnDone,
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

				assert.Equal(t, len(tt.expectedData.Tags), len(article.Tags))
				for i, tag := range tt.expectedData.Tags {
					assert.Equal(t, tag.ID, article.Tags[i].ID)
					assert.Equal(t, tag.Name, article.Tags[i].Name)
				}

				assert.Equal(t, tt.expectedData.Author.ID, article.Author.ID)
				assert.Equal(t, tt.expectedData.Author.Username, article.Author.Username)
				assert.Equal(t, tt.expectedData.Author.Email, article.Author.Email)
				assert.Equal(t, tt.expectedData.Author.Bio, article.Author.Bio)
				assert.Equal(t, tt.expectedData.Author.Image, article.Author.Image)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
