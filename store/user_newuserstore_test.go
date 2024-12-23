package store

import (
	"sync"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestNewUserStore(t *testing.T) {

	type testCase struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		scenario string
	}

	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	gormDB, err := gorm.Open("mysql", mockDB)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	tests := []testCase{
		{
			name:     "Scenario 1: Successfully Create New UserStore with Valid DB Connection",
			db:       gormDB,
			wantNil:  false,
			scenario: "Basic initialization with valid DB",
		},
		{
			name:     "Scenario 2: Create UserStore with Nil DB Connection",
			db:       nil,
			wantNil:  false,
			scenario: "Handling nil DB connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Starting:", tt.scenario)

			userStore := NewUserStore(tt.db)

			if tt.wantNil {
				assert.Nil(t, userStore, "UserStore should be nil")
			} else {
				assert.NotNil(t, userStore, "UserStore should not be nil")
				assert.Equal(t, tt.db, userStore.db, "DB reference should match")
			}

			t.Log("Completed:", tt.scenario)
		})
	}

	t.Run("Scenario 3: Verify DB Reference Integrity", func(t *testing.T) {
		t.Log("Testing DB reference integrity")
		userStore := NewUserStore(gormDB)
		assert.Equal(t, gormDB, userStore.db, "DB reference should maintain integrity")
	})

	t.Run("Scenario 4: Multiple UserStore Instances Independence", func(t *testing.T) {
		t.Log("Testing multiple instance independence")
		store1 := NewUserStore(gormDB)
		store2 := NewUserStore(gormDB)
		assert.NotSame(t, store1, store2, "Different instances should not be the same")
	})

	t.Run("Scenario 5: UserStore Creation with Configured DB Properties", func(t *testing.T) {
		t.Log("Testing DB configuration preservation")
		configuredDB := gormDB.LogMode(true)
		userStore := NewUserStore(configuredDB)
		assert.Equal(t, configuredDB, userStore.db, "DB configuration should be preserved")
	})

	t.Run("Scenario 7: Concurrent UserStore Creation", func(t *testing.T) {
		t.Log("Testing concurrent creation")
		var wg sync.WaitGroup
		stores := make([]*UserStore, 10)
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				stores[index] = NewUserStore(gormDB)
			}(i)
		}
		wg.Wait()

		for _, store := range stores {
			assert.NotNil(t, store, "Concurrent creation should succeed")
			assert.Equal(t, gormDB, store.db, "DB reference should be maintained")
		}
	})
}
