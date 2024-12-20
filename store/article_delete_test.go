package store

import (
	"testing"
	"time"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestDelete(t *testing.T) {
	type testCase struct {
		name    string
		article *model.Article
		setupFn func(mock sqlmock.Sqlmock)
		wantErr bool
	}

	tests := []testCase{
		{
			name: "Successfully Delete Existing Article",
			article: &model.Article{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Title:       "Test Article",
				Description: "Test Description",
				Body:        "Test Body",
				UserID:      1,
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `articles`").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Delete Non-existent Article",
			article: &model.Article{
				Model: gorm.Model{
					ID: 999,
				},
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `articles`").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Delete Article with Database Error",
			article: &model.Article{
				Model: gorm.Model{
					ID: 2,
				},
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `articles`").
					WithArgs(2).
					WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Delete Article with Associated Records",
			article: &model.Article{
				Model: gorm.Model{
					ID: 3,
				},
				Comments: []model.Comment{
					{Model: gorm.Model{ID: 1}},
				},
				Tags: []model.Tag{
					{Model: gorm.Model{ID: 1}},
				},
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `articles`").
					WithArgs(3).
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

			err = store.Delete(tt.article)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}
