package store

import (
	"database/sql"
	"errors"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreGetCommentByID(t *testing.T) {
	tests := []struct {
		name          string
		commentID     uint
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
		expectComment bool
	}{
		{
			name:      "Successfully retrieve existing comment",
			commentID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "body", "user_id", "article_id"}
				mock.ExpectQuery("^SELECT \\* FROM `comments` WHERE \\(`comments`\\.`id` = \\?\\) AND `comments`\\.`deleted_at` IS NULL LIMIT 1").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, time.Now(), time.Now(), nil, "Test comment", 1, 1))
			},
			expectedError: nil,
			expectComment: true,
		},
		{
			name:      "Attempt to retrieve non-existent comment",
			commentID: 999,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `comments` WHERE \\(`comments`\\.`id` = \\?\\) AND `comments`\\.`deleted_at` IS NULL LIMIT 1").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
			expectComment: false,
		},
		{
			name:      "Handle database connection error",
			commentID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `comments` WHERE \\(`comments`\\.`id` = \\?\\) AND `comments`\\.`deleted_at` IS NULL LIMIT 1").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: sql.ErrConnDone,
			expectComment: false,
		},
		{
			name:      "Retrieve comment with zero ID",
			commentID: 0,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `comments` WHERE \\(`comments`\\.`id` = \\?\\) AND `comments`\\.`deleted_at` IS NULL LIMIT 1").
					WithArgs(0).
					WillReturnError(errors.New("invalid comment ID"))
			},
			expectedError: errors.New("invalid comment ID"),
			expectComment: false,
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
			gormDB.LogMode(true)
			defer gormDB.Close()

			tt.setupMock(mock)

			store := &ArticleStore{
				db: gormDB,
			}

			comment, err := store.GetCommentByID(tt.commentID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.expectComment {
				assert.NotNil(t, comment)
				assert.IsType(t, &model.Comment{}, comment)
			} else {
				assert.Nil(t, comment)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
