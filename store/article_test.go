package store

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f


 */
func TestArticleStoreDeleteFavorite(t *testing.T) {

	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err, "failed to create sqlmock.New")

	db, gormErr := gorm.Open("postgres", sqlDB)
	assert.NoError(t, gormErr, "failed to open gorm with sqlite 3")

	var tests = []struct {
		name string

		mock func()

		expectError bool

		expectRollback bool

		setup func(store *ArticleStore, user *model.User, article *model.Article)
	}{
		{
			name: "Successful Deletion of a Favorite",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("^SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectExec("^DELETE FROM").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("^UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError:    false,
			expectRollback: false,
			setup: func(store *ArticleStore, user *model.User, article *model.Article) {
				article.ID = 1
				user.ID = 2
			},
		},
		{
			name: "Deletion of a Non-favorite Article",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("^SELECT").WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError:    true,
			expectRollback: true,
			setup: func(store *ArticleStore, user *model.User, article *model.Article) {
				article.ID = 3
				user.ID = 4
			},
		},
		{
			name: "Database Rollback on Error During Deletion",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("^SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectExec("^DELETE FROM").WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError:    true,
			expectRollback: true,
			setup: func(store *ArticleStore, user *model.User, article *model.Article) {
				article.ID = 5
				user.ID = 6
			},
		},
		{
			name: "Database Rollback During Updating of Article",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("^SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectExec("^DELETE FROM").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("^UPDATE").WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError:    true,
			expectRollback: true,
			setup: func(store *ArticleStore, user *model.User, article *model.Article) {
				article.ID = 7
				user.ID = 8
			},
		},
		{
			name: "Decrease of Favorited Count",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery("^SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "favorites_count"}).AddRow(1, 10))
				mock.ExpectExec("^DELETE FROM").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("^UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError:    false,
			expectRollback: false,
			setup: func(store *ArticleStore, user *model.User, article *model.Article) {
				article.ID = 9
				article.FavoritesCount = 10
				user.ID = 10
			},
		},
	}

	as := ArticleStore{db: db}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Log("Current test scenario: ", test.name)
			article := &model.Article{}
			user := &model.User{}

			test.setup(&as, user, article)
			test.mock()
			err := as.DeleteFavorite(article, user)
			if test.expectError {
				assert.Error(t, err, "expected to return an error")
			} else {
				assert.NoError(t, err, "expected to return without an error")
			}
			if test.expectRollback {

				assert.Error(t, mock.ExpectationsWereMet(), "expect Rollback to be called but did not")
			}
			if !test.expectError && test.name == "Decrease of Favorited Count" {
				assert.Equal(t, int32(9), article.FavoritesCount, "expected favorites count to decrease by one")
			}
		})
	}
}

