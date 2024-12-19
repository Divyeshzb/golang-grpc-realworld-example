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
			errorMsg:    "empty TOML file",
		},
		{
			name: "Invalid TOML Format",
			tomlContent: `[[Users]]
				invalid_toml_syntax`,
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectError: true,
			errorMsg:    "TOML parsing error",
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
			errorMsg:    "database error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			tmpfile, err := ioutil.TempFile("", "users.*.toml")
			if err != nil {
				t.Fatal("Failed to create temp file:", err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tc.tomlContent)); err != nil {
				t.Fatal("Failed to write to temp file:", err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal("Failed to close temp file:", err)
			}

			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer mockDB.Close()

			gormDB, err := gorm.Open("mysql", mockDB)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.setupMock(mock)

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
