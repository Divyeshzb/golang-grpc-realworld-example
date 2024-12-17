package store

import (
	"database/sql"
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestArticleStoreGetByID(t *testing.T) {

	type testCase struct {
		name         string
		id           uint
		mockSetup    func(sqlmock.Sqlmock)
		expectedErr  error
		expectedData *model.Article
	}

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

	store := &ArticleStore{db: gormDB}

	tests := []testCase{
		{
			name: "Successful Article Retrieval",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery("SELECT \\* FROM `articles` WHERE").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"title", "body"}).
						AddRow("Test Article", "Test Content"))

				mock.ExpectQuery("SELECT \\* FROM `tags`").
					WillReturnRows(sqlmock.NewRows([]string{"name"}).
						AddRow("golang"))

				mock.ExpectQuery("SELECT \\* FROM `users`").
					WillReturnRows(sqlmock.NewRows([]string{"username"}).
						AddRow("testuser"))
			},
			expectedData: &model.Article{
				Title: "Test Article",
				Body:  "Test Content",
				Tags: []model.Tag{
					{Name: "golang"},
				},
				Author: model.User{
					Username: "testuser",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Non-existent Article",
			id:   99999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `articles` WHERE").
					WithArgs(99999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedData: nil,
			expectedErr:  gorm.ErrRecordNotFound,
		},
		{
			name: "Database Error",
			id:   2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `articles` WHERE").
					WithArgs(2).
					WillReturnError(sql.ErrConnDone)
			},
			expectedData: nil,
			expectedErr:  sql.ErrConnDone,
		},
		{
			name: "Zero ID Input",
			id:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `articles` WHERE").
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedData: nil,
			expectedErr:  gorm.ErrRecordNotFound,
		},
		{
			name: "Article with Empty Relationships",
			id:   3,
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery("SELECT \\* FROM `articles` WHERE").
					WithArgs(3).
					WillReturnRows(sqlmock.NewRows([]string{"title", "body"}).
						AddRow("Empty Article", "No relationships"))

				mock.ExpectQuery("SELECT \\* FROM `tags`").
					WillReturnRows(sqlmock.NewRows([]string{"name"}))

				mock.ExpectQuery("SELECT \\* FROM `users`").
					WillReturnRows(sqlmock.NewRows([]string{"username"}))
			},
			expectedData: &model.Article{
				Title: "Empty Article",
				Body:  "No relationships",
				Tags:  []model.Tag{},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			tc.mockSetup(mock)

			article, err := store.GetByID(tc.id)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
				assert.Nil(t, article)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, article)
				assert.Equal(t, tc.expectedData.Title, article.Title)
				assert.Equal(t, tc.expectedData.Body, article.Body)
				assert.Equal(t, len(tc.expectedData.Tags), len(article.Tags))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
