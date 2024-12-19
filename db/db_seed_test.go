package db

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestSeed(t *testing.T) {

	type testCase struct {
		name        string
		tomlContent string
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}

	tests := []testCase{
		{
			name: "Successful Database Seeding",
			tomlContent: `[[Users]]
				id = "550e8400-e29b-41d4-a716-446655440000"
				email = "test@example.com"
				username = "testuser"
				password = "password123"`,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name:        "Empty TOML File",
			tomlContent: "",
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectError: true,
			errorMsg:    "toml: cannot decode empty string",
		},
		{
			name: "Invalid TOML Format",
			tomlContent: `[[Users]]
				invalid_toml_format`,
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectError: true,
			errorMsg:    "toml:",
		},
		{
			name: "Database Error",
			tomlContent: `[[Users]]
				id = "550e8400-e29b-41d4-a716-446655440000"
				email = "test@example.com"
				username = "testuser"
				password = "password123"`,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").WillReturnError(sqlmock.ErrCancelled)
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "sql: transaction has been cancelled",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			tmpfile, err := ioutil.TempFile("", "users.*.toml")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if err := ioutil.WriteFile(tmpfile.Name(), []byte(tc.tomlContent), 0644); err != nil {
				t.Fatal(err)
			}

			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database: %v", err)
			}
			defer sqlDB.Close()

			tc.setupMock(mock)

			gormDB, err := gorm.Open("mysql", sqlDB)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			originalPath := "db/seed/users.toml"
			os.Rename("db/seed/users.toml", "db/seed/users.toml.bak")
			os.Symlink(tmpfile.Name(), originalPath)
			defer func() {
				os.Remove(originalPath)
				os.Rename("db/seed/users.toml.bak", originalPath)
			}()

			err = Seed(gormDB)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
