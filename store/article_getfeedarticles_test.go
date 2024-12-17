package store

import (
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreGetFeedArticles(t *testing.T) {

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
		name      string
		userIDs   []uint
		limit     int64
		offset    int64
		mockSetup func(sqlmock.Sqlmock)
		expected  []model.Article
		wantErr   bool
	}{
		{
			name:    "Successfully retrieve feed articles",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				}).AddRow(
					1, time.Now(), time.Now(), nil,
					"Test Article", "Test Description", "Test Body", 1, 0,
				)

				mock.ExpectQuery("SELECT \\* FROM `articles`").
					WithArgs(1, 2).
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"username", "email", "password", "bio", "image",
				}).AddRow(
					1, time.Now(), time.Now(), nil,
					"testuser", "test@example.com", "hashedpass", "bio", "image.jpg",
				)
				mock.ExpectQuery("SELECT \\* FROM `users`").
					WillReturnRows(authorRows)
			},
			expected: []model.Article{
				{
					Title:       "Test Article",
					Description: "Test Description",
					Body:        "Test Body",
					UserID:      1,
					Author: model.User{
						Username: "testuser",
						Email:    "test@example.com",
						Bio:      "bio",
						Image:    "image.jpg",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Empty result when no articles exist",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})
				mock.ExpectQuery("SELECT \\* FROM `articles`").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expected: []model.Article{},
			wantErr:  false,
		},
		{
			name:    "Database error handling",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `articles`").
					WithArgs(1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name:    "Empty userIDs slice handling",
			userIDs: []uint{},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})
				mock.ExpectQuery("SELECT \\* FROM `articles`").
					WillReturnRows(rows)
			},
			expected: []model.Article{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup(mock)

			got, err := store.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expected), len(got))

				if len(tt.expected) > 0 {
					assert.Equal(t, tt.expected[0].Title, got[0].Title)
					assert.Equal(t, tt.expected[0].Description, got[0].Description)
					assert.Equal(t, tt.expected[0].Body, got[0].Body)
					assert.Equal(t, tt.expected[0].UserID, got[0].UserID)

					assert.Equal(t, tt.expected[0].Author.Username, got[0].Author.Username)
					assert.Equal(t, tt.expected[0].Author.Email, got[0].Author.Email)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
