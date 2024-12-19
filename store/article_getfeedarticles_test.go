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

				mock.ExpectQuery("SELECT").
					WithArgs(1, 2).
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{
					"id", "username", "email", "password", "bio", "image",
				}).AddRow(1, "testuser", "test@test.com", "password", "bio", "image")

				mock.ExpectQuery("SELECT").
					WillReturnRows(authorRows)
			},
			expected: []model.Article{
				{
					Model: gorm.Model{ID: 1},
					Title: "Test Article",
					Author: model.User{
						Model:    gorm.Model{ID: 1},
						Username: "testuser",
					},
					UserID: 1,
				},
			},
			wantErr: false,
		},
		{
			name:    "Empty userIDs slice",
			userIDs: []uint{},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expected: []model.Article{},
			wantErr:  false,
		},
		{
			name:    "Database error",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(gorm.ErrRecordNotFound)
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name:    "Pagination test",
			userIDs: []uint{1},
			limit:   5,
			offset:  5,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id", "favorites_count",
				})
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
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
					assert.Equal(t, tt.expected[0].UserID, got[0].UserID)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed. Got %d articles", tt.name, len(got))
		})
	}
}
