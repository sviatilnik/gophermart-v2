package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/sviatilnik/gophermart/internal/domain/user"
	"testing"
)

func TestPostgresUserRepository_Exists(t *testing.T) {
	type args struct {
		ctx   context.Context
		login string
	}
	tests := []struct {
		name      string
		args      args
		want      bool
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
		err       error
	}{
		{
			name: "success",
			args: args{
				ctx:   context.TODO(),
				login: "Тест",
			},
			want: true,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE login = (.+)$").
					WithArgs("Тест").
					WillReturnRows(mock.NewRows([]string{"id", "login", "password"}).
						AddRow(1, "Тест", ""))
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "not found",
			args: args{
				ctx:   context.TODO(),
				login: "Тест",
			},
			want: false,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE login = (.+)$").
					WithArgs(sqlmock.AnyArg()).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			err:     user.ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			r := NewPostgresUserRepository(db)
			exists, er := r.Exists(tt.args.ctx, user.Login(tt.args.login))
			if tt.wantErr {
				assert.ErrorIs(t, er, tt.err)
			} else {
				assert.NoError(t, er)
				assert.Equal(t, tt.want, exists)
			}
		})
	}
}

func TestPostgresUserRepository_FindByID(t *testing.T) {
	password, _ := user.NewPassword("test")

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name      string
		args      args
		want      *user.User
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
		err       error
	}{
		{
			name: "success",
			args: args{
				ctx: context.TODO(),
				id:  "1",
			},
			want: &user.User{
				ID:       "1",
				Login:    "Тест",
				Password: password,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE id = (.+)$").
					WithArgs("1").
					WillReturnRows(mock.NewRows([]string{"id", "login", "password"}).
						AddRow(1, "Тест", password))
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "not found",
			args: args{
				ctx: context.TODO(),
				id:  "2",
			},
			want: nil,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE id = (.+)$").
					WithArgs("2").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			err:     user.ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			r := NewPostgresUserRepository(db)
			usr, er := r.FindByID(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				assert.ErrorIs(t, er, tt.err)
			} else {
				assert.NoError(t, er)
				assert.Equal(t, tt.want, usr)
			}
		})
	}
}

func TestPostgresUserRepository_FindByLogin(t *testing.T) {
	password, _ := user.NewPassword("test")

	type args struct {
		ctx   context.Context
		login user.Login
	}

	tests := []struct {
		name      string
		args      args
		want      *user.User
		mockSetup func(mock sqlmock.Sqlmock)
		wantErr   bool
		err       error
	}{
		{
			name: "success",
			args: args{
				ctx:   context.TODO(),
				login: user.Login("Тест"),
			},
			want: &user.User{
				ID:       "1",
				Login:    "Тест",
				Password: password,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE login = (.+)$").
					WithArgs("Тест").
					WillReturnRows(mock.NewRows([]string{"id", "login", "password"}).
						AddRow(1, "Тест", password))
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "not found",
			args: args{
				ctx:   context.TODO(),
				login: user.Login("Тест"),
			},
			want: nil,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE login = (.+)$").
					WithArgs("Тест").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			err:     user.ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			r := NewPostgresUserRepository(db)
			usr, er := r.FindByLogin(tt.args.ctx, tt.args.login)
			if tt.wantErr {
				assert.ErrorIs(t, er, tt.err)
			} else {
				assert.NoError(t, er)
				assert.Equal(t, tt.want, usr)
			}
		})
	}
}

func TestPostgresUserRepository_Save(t *testing.T) {
	type args struct {
		ctx  context.Context
		user *user.User
	}
	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		args      args
		wantErr   bool
		err       error
	}{
		{
			name: "success",
			args: args{
				ctx:  context.TODO(),
				user: &user.User{ID: "5", Login: "Test", Password: "Text"},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^INSERT INTO users (.*) VALUES (.*)$").
					WithArgs("5", "Test", "Text").
					WillReturnResult(sqlmock.NewResult(5, 1))
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "error",
			args: args{
				ctx:  context.TODO(),
				user: &user.User{ID: "5", Login: "Test", Password: "Text"},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^INSERT INTO users (.*) VALUES (.*)$").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("insert error"))
			},
			wantErr: true,
			err:     errors.New("insert error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			r := NewPostgresUserRepository(db)

			err = r.Save(tt.args.ctx, tt.args.user)
			if tt.wantErr {
				assert.Equal(t, err.Error(), tt.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPostgresUserRepository_Delete(t *testing.T) {
	type args struct {
		ctx  context.Context
		user *user.User
	}
	tests := []struct {
		name      string
		mockSetup func(mock sqlmock.Sqlmock)
		args      args
		wantErr   bool
		err       error
	}{
		{
			name: "success",
			args: args{
				ctx:  context.TODO(),
				user: &user.User{ID: "5", Login: "Test", Password: "Text"},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^DELETE FROM users WHERE id = (.*)$").
					WithArgs("5").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "error",
			args: args{
				ctx:  context.TODO(),
				user: &user.User{ID: "5", Login: "Test", Password: "Text"},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^DELETE FROM users WHERE id = (.*)$").
					WithArgs(sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
			err:     user.ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			r := NewPostgresUserRepository(db)

			err = r.Delete(tt.args.ctx, tt.args.user.ID)
			if tt.wantErr {
				assert.Equal(t, err.Error(), tt.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
