package store

import (
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name        string
		article     *model.Article
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "Successful Update",
			article: &model.Article{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Title:  "Updated Title",
				Body:   "Updated content",
				Author: model.User{},
				Tags:   []model.Tag{{Name: "test"}},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Update Non-Existent Article",
			article: &model.Article{
				Model: gorm.Model{
					ID: 999,
				},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "Update with Invalid Data",
			article: &model.Article{
				Model: gorm.Model{
					ID: 1,
				},
				Title: "",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "Update with No Changes",
			article: &model.Article{
				Model: gorm.Model{
					ID: 1,
				},
				Title:  "Existing Title",
				Body:   "Existing content",
				Author: model.User{},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			err := store.Update(tt.article)

			if tt.expectError {
				assert.Error(t, err)
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				t.Logf("Update successful for article ID: %d", tt.article.ID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
