package store

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"testing"
)

type expectedResults struct {
	err               error
	isUserInFavorites bool
	favoritesCount    int
}
/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f


 */
func TestDeleteFavorite(t *testing.T) {

	scenarios := []struct {
		desc           string
		prepareMock    func(mock sqlmock.Sqlmock)
		expectedResult expectedResults
	}{
		{

			desc: "Successfully deleting a favorite",
			prepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM 'favorited_users' WHERE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("^UPDATE 'articles' SET").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedResult: expectedResults{
				err:               nil,
				isUserInFavorites: false,
				favoritesCount:    0,
			},
		},
	}

	a := model.Article{
		Model:          gorm.Model{ID: 1},
		FavoritedUsers: []model.User{{Model: gorm.Model{ID: 1}}},
		FavoritesCount: 1,
	}

	u := model.User{
		Model: gorm.Model{ID: 1},
	}

	for _, s := range scenarios {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		gdb, err := gorm.Open("postgres", db)
		if err != nil {
			t.Errorf("an error '%s' was not expected when opening gorm database", err)
		}

		s.prepareMock(mock)

		store := ArticleStore{db: gdb}

		err = store.DeleteFavorite(&a, &u)

		assert.Equal(t, s.expectedResult.err, err)
	}
}

