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

func TestArticleStoreDeleteComment(t *testing.T) {

	tests := []struct {
		name    string
		comment *model.Comment
		setupFn func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Successfully Delete Existing Comment",
			comment: &model.Comment{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Body:      "Test Comment",
				UserID:    1,
				ArticleID: 1,
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `comments` SET").
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Delete Non-existent Comment",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 999,
				},
				Body:      "Non-existent Comment",
				UserID:    1,
				ArticleID: 1,
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `comments` SET").
					WithArgs(sqlmock.AnyArg(), 999).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Delete Comment with NULL Values",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 2,
				},
				Body:      "Minimal Comment",
				UserID:    1,
				ArticleID: 1,
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `comments` SET").
					WithArgs(sqlmock.AnyArg(), 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Delete Comment with Invalid DB Connection",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 3,
				},
				Body:      "Test Comment",
				UserID:    1,
				ArticleID: 1,
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `comments` SET").
					WithArgs(sqlmock.AnyArg(), 3).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Delete Comment with Associated Relations",
			comment: &model.Comment{
				Model: gorm.Model{
					ID: 4,
				},
				Body:      "Comment with Relations",
				UserID:    1,
				ArticleID: 1,
				Author:    model.User{},
				Article:   model.Article{},
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `comments` SET").
					WithArgs(sqlmock.AnyArg(), 4).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
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
				t.Fatalf("Failed to open GORM DB: %v", err)
			}
			defer gormDB.Close()

			tt.setupFn(mock)

			store := &ArticleStore{
				db: gormDB,
			}

			err = store.DeleteComment(tt.comment)

			if tt.wantErr {
				assert.Error(t, err)
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				t.Log("Comment deleted successfully")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
