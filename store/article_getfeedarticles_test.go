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

func TestGetFeedArticles(t *testing.T) {
	type testCase struct {
		name      string
		userIDs   []uint
		limit     int64
		offset    int64
		mockSetup func(sqlmock.Sqlmock)
		expected  []model.Article
		expectErr bool
	}

	timeNow := time.Now()

	tests := []testCase{
		{
			name:    "Successfully retrieve feed articles",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				articleRows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id",
					"favorites_count",
				}).
					AddRow(1, timeNow, timeNow, nil, "Test Article 1", "Description 1", "Body 1", 1, 0).
					AddRow(2, timeNow, timeNow, nil, "Test Article 2", "Description 2", "Body 2", 2, 0)

				authorRows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"username", "email", "bio", "image",
				}).
					AddRow(1, timeNow, timeNow, nil, "user1", "user1@example.com", "", "").
					AddRow(2, timeNow, timeNow, nil, "user2", "user2@example.com", "", "")

				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE \(user_id in \(\?,\?\)\) LIMIT \? OFFSET \?`).
					WithArgs(1, 2, 10, 0).
					WillReturnRows(articleRows)

				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" IN \(\?,\?\)`).
					WithArgs(1, 2).
					WillReturnRows(authorRows)
			},
			expected: []model.Article{
				{
					Model:       gorm.Model{ID: 1, CreatedAt: timeNow, UpdatedAt: timeNow},
					Title:       "Test Article 1",
					Description: "Description 1",
					Body:        "Body 1",
					UserID:      1,
					Author: model.User{
						Model:    gorm.Model{ID: 1, CreatedAt: timeNow, UpdatedAt: timeNow},
						Username: "user1",
						Email:    "user1@example.com",
					},
				},
				{
					Model:       gorm.Model{ID: 2, CreatedAt: timeNow, UpdatedAt: timeNow},
					Title:       "Test Article 2",
					Description: "Description 2",
					Body:        "Body 2",
					UserID:      2,
					Author: model.User{
						Model:    gorm.Model{ID: 2, CreatedAt: timeNow, UpdatedAt: timeNow},
						Username: "user2",
						Email:    "user2@example.com",
					},
				},
			},
			expectErr: false,
		},
		{
			name:    "Empty result set",
			userIDs: []uint{99},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id",
					"favorites_count",
				})
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE \(user_id in \(\?\)\) LIMIT \? OFFSET \?`).
					WithArgs(99, 10, 0).
					WillReturnRows(rows)
			},
			expected:  []model.Article{},
			expectErr: false,
		},
		{
			name:    "Database error",
			userIDs: []uint{1},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE \(user_id in \(\?\)\) LIMIT \? OFFSET \?`).
					WithArgs(1, 10, 0).
					WillReturnError(sql.ErrConnDone)
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name:    "Empty userIDs array",
			userIDs: []uint{},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"title", "description", "body", "user_id",
					"favorites_count",
				})
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE \(user_id in \(\)\) LIMIT \? OFFSET \?`).
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			expected:  []model.Article{},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
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

			tc.mockSetup(mock)

			store := &ArticleStore{
				db: gormDB,
			}

			articles, err := store.GetFeedArticles(tc.userIDs, tc.limit, tc.offset)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, articles)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
