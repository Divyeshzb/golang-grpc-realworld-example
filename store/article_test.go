package store

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f


 */
func (m *MockDB) Association(column string) *gorm.Association {
	args := m.Called(column)
	return args.Get(0).(*gorm.Association)
}

func (m *MockDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Commit() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func TestDeleteFavorite(t *testing.T) {
	mockDB := new(MockDB)
	mockDB.DB = new(gorm.DB)
	s := &ArticleStore{db: mockDB.DB}

	t.Run("Delete favorite successfully", func(t *testing.T) {
		mockDB.On("Begin").Return(mockDB.DB)
		mockDB.On("Model", a).Return(mockDB.DB).Once()
		mockDB.On("Association", "FavoritedUsers").Return(mockDB.DB).Once()
		mockDB.On("Update", "favorites_count", gorm.Expr("favorites_count - ?", 1)).Return(mockDB.DB).Once()
		mockDB.On("Commit").Return(mockDB.DB)

		err := s.DeleteFavorite(a, u)
		assert.NoError(t, err)
		assert.Equal(t, a.FavoritesCount, 1)
	})

	t.Run("Error in deleting favorite", func(t *testing.T) {
		mockDB.On("Begin").Return(mockDB.DB)
		mockDB.On("Model", a).Return(mockDB.DB).Once()
		mockDB.On("Association", "FavoritedUsers").Return(errors.New("delete error").Error())
		mockDB.On("Rollback").Return(mockDB.DB)

		err := s.DeleteFavorite(a, u)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "delete error")
	})

	t.Run("Error in updating favorites count", func(t *testing.T) {
		mockDB.On("Begin").Return(mockDB.DB)
		mockDB.On("Model", a).Return(mockDB.DB).Once()
		mockDB.On("Association", "FavoritedUsers").Return(mockDB.DB).Once()
		mockDB.On("Update", "favorites_count", gorm.Expr("favorites_count - ?", 1)).Return(errors.New("update error"))
		mockDB.On("Rollback").Return(mockDB.DB)

		err := s.DeleteFavorite(a, u)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "update error")
	})

	t.Run("Unexpected increase in Favorites Count", func(t *testing.T) {
		a.FavoritesCount = -2

		mockDB.On("Begin").Return(mockDB.DB)
		mockDB.On("Model", a).Return(mockDB.DB).Once()
		mockDB.On("Association", "FavoritedUsers").Return(mockDB.DB).Once()
		mockDB.On("Update", "favorites_count", gorm.Expr("favorites_count - ?", 1)).Return(mockDB.DB).Once()
		mockDB.On("Commit").Return(mockDB.DB)

		err := s.DeleteFavorite(a, u)
		assert.NoError(t, err)
		assert.Equal(t, a.FavoritesCount, -3)
	})
}

func (m *MockDB) Update(column string, value interface{}) *gorm.DB {
	args := m.Called(column, value)
	return args.Get(0).(*gorm.DB)
}

