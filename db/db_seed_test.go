package db

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"os"
	"testing"
)

func TestSeed(t *testing.T) {

	type testCase struct {
		name        string
		tomlContent string
		setupDB     func(mock sqlmock.Sqlmock)
		wantErr     bool
		errMsg      string
	}

	tests := []testCase{
		{
			name: "Successful Database Seeding",
			tomlContent: `[[Users]]
				id = "550e8400-e29b-41d4-a716-446655440000"
				username = "testuser1"
				email = "test1@example.com"
				password = "hashedpassword1"
				bio = "test bio 1"
				image = "image1.jpg"
				created_at = 2024-01-01T00:00:00Z
				updated_at = 2024-01-01T00:00:00Z

				[[Users]]
				id = "550e8400-e29b-41d4-a716-446655440001"
				username = "testuser2"
				email = "test2@example.com"
				password = "hashedpassword2"
				bio = "test bio 2"
				image = "image2.jpg"
				created_at = 2024-01-01T00:00:00Z
				updated_at = 2024-01-01T00:00:00Z`,
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO `users`").WillReturnResult(sqlmock.NewResult(2, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:        "Empty TOML File",
			tomlContent: "",
			setupDB:     func(mock sqlmock.Sqlmock) {},
			wantErr:     true,
			errMsg:      "empty TOML file",
		},
		{
			name: "Invalid TOML Format",
			tomlContent: `[[Users]]
				invalid toml content`,
			setupDB: func(mock sqlmock.Sqlmock) {},
			wantErr: true,
			errMsg:  "TOML parsing error",
		},
		{
			name: "Database Error",
			tomlContent: `[[Users]]
				id = "550e8400-e29b-41d4-a716-446655440000"
				username = "testuser1"
				email = "test1@example.com"
				password = "hashedpassword1"`,
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").WillReturnError(sqlmock.ErrCancelled)
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tmpfile, err := ioutil.TempFile("", "users.*.toml")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.tomlContent)); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatalf("Failed to close temp file: %v", err)
			}

			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer mockDB.Close()

			db, err := gorm.Open("mysql", mockDB)
			if err != nil {
				t.Fatalf("Failed to open GORM DB: %v", err)
			}
			defer db.Close()

			tt.setupDB(mock)

			err = Seed(db)

			if (err != nil) != tt.wantErr {
				t.Errorf("Seed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case '%s' completed successfully", tt.name)
		})
	}
}
