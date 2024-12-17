package store

import (
	"testing"
	"github.com/jinzhu/gorm"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewArticleStore(t *testing.T) {

	tests := []struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		scenario string
	}{
		{
			name:     "Scenario 1: Successfully Create New ArticleStore with Valid DB Connection",
			db:       setupTestDB(t),
			wantNil:  false,
			scenario: "Valid DB connection should create a proper ArticleStore instance",
		},
		{
			name:     "Scenario 2: Create ArticleStore with Nil DB Parameter",
			db:       nil,
			wantNil:  false,
			scenario: "Nil DB should still create an ArticleStore instance but with nil db field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Starting:", tt.scenario)

			store := NewArticleStore(tt.db)

			if tt.wantNil {
				assert.Nil(t, store, "Expected nil ArticleStore")
			} else {
				assert.NotNil(t, store, "Expected non-nil ArticleStore")
				assert.Equal(t, tt.db, store.db, "DB connection mismatch")
			}

			if tt.db != nil {
				store1 := NewArticleStore(tt.db)
				store2 := NewArticleStore(tt.db)
				assert.NotEqual(t, store1, store2, "Different instances should have different memory addresses")
				assert.Equal(t, store1.db, store2.db, "DB connections should be the same")
			}

			if tt.db != nil {
				assert.Equal(t, tt.db, store.db, "DB connection should persist unchanged")

			}

			func() {
				var stores []*ArticleStore
				for i := 0; i < 100; i++ {
					stores = append(stores, NewArticleStore(tt.db))
				}

				stores = nil
			}()

			t.Log("Successfully completed:", tt.scenario)
		})
	}
}
func setupTestDB(t *testing.T) *gorm.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled mock expectations: %v", err)
		}
	})

	return gormDB
}
