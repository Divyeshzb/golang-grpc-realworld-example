package store

import (
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/DATA-DOG/go-sqlmock"
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
			name: "Scenario 1: Successfully Create New ArticleStore with Valid DB Connection",
			db: func() *gorm.DB {
				db, _, err := sqlmock.New()
				if err != nil {
					t.Fatalf("Failed to create mock DB: %v", err)
				}
				gormDB, err := gorm.Open("mysql", db)
				if err != nil {
					t.Fatalf("Failed to create GORM DB: %v", err)
				}
				return gormDB
			}(),
			wantNil:  false,
			scenario: "Valid DB connection should create valid ArticleStore",
		},
		{
			name:     "Scenario 2: Create ArticleStore with Nil DB Connection",
			db:       nil,
			wantNil:  false,
			scenario: "Nil DB should still create ArticleStore instance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Logf("Testing scenario: %s", tt.scenario)

			got := NewArticleStore(tt.db)

			if tt.wantNil {
				assert.Nil(t, got, "Expected nil ArticleStore")
			} else {
				assert.NotNil(t, got, "Expected non-nil ArticleStore")
				assert.Equal(t, tt.db, got.db, "DB reference mismatch")
			}

			if tt.db != nil {
				assert.Same(t, tt.db, got.db, "DB reference should be the same instance")
			}

			t.Logf("Successfully completed test: %s", tt.name)
		})
	}

	t.Run("Scenario 3: Verify DB Reference Integrity", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		gormDB, err := gorm.Open("mysql", db)
		if err != nil {
			t.Fatalf("Failed to create GORM DB: %v", err)
		}

		store := NewArticleStore(gormDB)
		assert.Equal(t, gormDB, store.db, "DB reference should match original")
	})

	t.Run("Scenario 4: Multiple ArticleStore Instances Independence", func(t *testing.T) {
		db1, _, _ := sqlmock.New()
		db2, _, _ := sqlmock.New()
		gormDB1, _ := gorm.Open("mysql", db1)
		gormDB2, _ := gorm.Open("mysql", db2)

		store1 := NewArticleStore(gormDB1)
		store2 := NewArticleStore(gormDB2)

		assert.NotEqual(t, store1.db, store2.db, "Different stores should have different DB instances")
	})

	t.Run("Scenario 5: ArticleStore Creation with Configured DB", func(t *testing.T) {
		db, _, _ := sqlmock.New()
		gormDB, _ := gorm.Open("mysql", db)
		gormDB.LogMode(true)

		store := NewArticleStore(gormDB)
		assert.Equal(t, gormDB.LogMode(false), store.db.LogMode(false), "DB configuration should be preserved")
	})

}
