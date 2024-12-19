package store

import (
	"database/sql"
	"regexp"
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreCreate(t *testing.T) {

	tests := []struct {
		name    string
		user    *model.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful User Creation",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			wantErr: false,
		},
		{
			name: "Duplicate Username",
			user: &model.User{
				Username: "existinguser",
				Email:    "new@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			wantErr: true,
			errMsg:  "duplicate key value violates unique constraint",
		},
		{
			name: "Duplicate Email",
			user: &model.User{
				Username: "newuser",
				Email:    "existing@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			wantErr: true,
			errMsg:  "duplicate key value violates unique constraint",
		},
		{
			name: "Missing Required Fields",
			user: &model.User{
				Username: "",
				Email:    "",
				Password: "",
			},
			wantErr: true,
			errMsg:  "not null constraint violation",
		},
		{
			name: "Maximum Field Lengths",
			user: &model.User{
				Username: "maxlengthusername123456789",
				Email:    "verylongemail@verylongdomain.com",
				Password: "verylongpassword123456789",
				Bio:      "A very long bio text that tests the maximum field length",
				Image:    "very-long-image-url-path/image.jpg",
			},
			wantErr: false,
		},
		{
			name: "Special Characters",
			user: &model.User{
				Username: "user@#$%",
				Email:    "special.chars+test@domain.com",
				Password: "pass!@#$%^&*()",
				Bio:      "Bio with émojis 🎉 and ünicode",
				Image:    "image-with-späces.jpg",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			store := &UserStore{db: gormDB}

			if !tt.wantErr {

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						tt.user.Username, tt.user.Email, tt.user.Password,
						tt.user.Bio, tt.user.Image).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			} else {

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			}

			err = store.Create(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				t.Log("User created successfully")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
