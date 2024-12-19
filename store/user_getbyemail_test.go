package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStoreGetByEmail(t *testing.T) {
	type args struct {
		email string
	}

	tests := []struct {
		name    string
		args    args
		mock    func(sqlmock.Sqlmock)
		want    *model.User
		wantErr error
	}{
		{
			name: "Success - Valid Email",
			args: args{
				email: "test@example.com",
			},
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email"}).
					AddRow(1, "test@example.com")
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			want: &model.User{
				Model: gorm.Model{ID: 1},
				Email: "test@example.com",
			},
			wantErr: nil,
		},
		{
			name: "Error - User Not Found",
			args: args{
				email: "nonexistent@example.com",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE").
					WithArgs("nonexistent@example.com").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Error - Database Error",
			args: args{
				email: "test@example.com",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE").
					WithArgs("test@example.com").
					WillReturnError(errors.New("database connection error"))
			},
			want:    nil,
			wantErr: errors.New("database connection error"),
		},
		{
			name: "Error - Empty Email",
			args: args{
				email: "",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE").
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Success - Special Characters in Email",
			args: args{
				email: "test+special@example.com",
			},
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email"}).
					AddRow(2, "test+special@example.com")
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE").
					WithArgs("test+special@example.com").
					WillReturnRows(rows)
			},
			want: &model.User{
				Model: gorm.Model{ID: 2},
				Email: "test+special@example.com",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			s := &UserStore{
				db: gormDB,
			}

			if tt.mock != nil {
				tt.mock(mock)
			}

			got, err := s.GetByEmail(tt.args.email)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}

			t.Logf("Test case '%s' completed", tt.name)
			if err != nil {
				t.Logf("Error occurred: %v", err)
			}
		})
	}
}
