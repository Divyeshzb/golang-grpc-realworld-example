package store

import (
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"reflect"
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
			name:     "Scenario 1: Successfully Create New UserStore",
			db:       gormDB,
			wantNil:  false,
			scenario: "Valid DB Connection",
		},
		{
			name:     "Scenario 2: Create UserStore with Nil DB",
			db:       nil,
			wantNil:  false,
			scenario: "Nil DB Parameter",
		},
		{
			name:     "Scenario 3: Verify DB Reference Integrity",
			db:       gormDB,
			wantNil:  false,
			scenario: "DB Reference Check",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Starting:", tc.scenario)

			userStore := NewUserStore(tc.db)

			if tc.wantNil {
				assert.Nil(t, userStore, "UserStore should be nil")
			} else {
				assert.NotNil(t, userStore, "UserStore should not be nil")

				if tc.db != nil {
					assert.Equal(t, tc.db, userStore.db, "DB reference mismatch")
				}
			}

			switch tc.scenario {
			case "Valid DB Connection":
				assert.NotNil(t, userStore.db, "DB connection should be initialized")

			case "Nil DB Parameter":

				assert.Nil(t, userStore.db, "DB should be nil")

			case "DB Reference Check":

				assert.True(t, reflect.DeepEqual(tc.db, userStore.db), "DB reference should match exactly")
			}

			t.Log("Completed:", tc.scenario)
		})
	}

	t.Run("Scenario 4: Multiple UserStore Instances Independence", func(t *testing.T) {
		t.Log("Testing multiple UserStore instances")

		mockDB1, _, _ := sqlmock.New()
		mockDB2, _, _ := sqlmock.New()
		defer mockDB1.Close()
		defer mockDB2.Close()

		gormDB1, _ := gorm.Open("mysql", mockDB1)
		gormDB2, _ := gorm.Open("mysql", mockDB2)
		defer gormDB1.Close()
		defer gormDB2.Close()

		store1 := NewUserStore(gormDB1)
		store2 := NewUserStore(gormDB2)

		assert.NotEqual(t, store1.db, store2.db, "Different UserStore instances should have independent DB references")
		t.Log("Multiple instances test completed")
	})

	t.Run("Scenario 5: UserStore Creation with Configured DB", func(t *testing.T) {
		t.Log("Testing DB configuration persistence")

		configuredDB := gormDB
		configuredDB.LogMode(true)

		userStore := NewUserStore(configuredDB)

		assert.Equal(t, configuredDB.LogMode(true), userStore.db.LogMode(true),
			"DB configuration should be maintained")
		t.Log("Configuration persistence test completed")
	})

	t.Run("Scenario 6: Memory Resource Management", func(t *testing.T) {
		t.Log("Testing resource management")

		for i := 0; i < 100; i++ {
			mockDB, _, _ := sqlmock.New()
			gormDB, _ := gorm.Open("mysql", mockDB)
			store := NewUserStore(gormDB)

			mockDB.Close()
			gormDB.Close()

			assert.NotNil(t, store, "UserStore should be created successfully")
		}
		t.Log("Resource management test completed")
	})
}
