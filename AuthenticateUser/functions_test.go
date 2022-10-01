package main

import (
	"database/sql"
	"os"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
)

func TestConnectDatabase(t *testing.T) {
	os.Setenv(`DB_USER`, `root`)
	os.Setenv(`DB_PASSWORD`, `1234`)
	os.Setenv(`DB_HOST`, `dbhost`)
	os.Setenv(`DB_NAME`, `dbname`)
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Error test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ConnectDatabase()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConnectDatabase error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestGetUserData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{`document`, `name`, `role`}
	columns_error := []string{`document`, `name`, `role`, `some`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`, `some_pass`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`123456`, `some_name`, `role`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`1234`, `some_pass`).
		WillReturnRows(sqlmock.NewRows(columns_error).AddRow(`123456`, `some_name`, `role`, `some`))

	type args struct {
		conn *sql.DB
		data RequestBody
	}
	tests := []struct {
		name      string
		args      args
		wantFound bool
		wantUser  User
		wantErr   bool
	}{
		{
			name: "Success test - user found",
			args: args{
				conn: db,
				data: RequestBody{
					Document: `123456`,
					Password: `some_pass`,
				},
			},
			wantFound: true,
			wantUser: User{
				Name:   `some_name`,
				UserId: `123456`,
				Role:   `role`,
			},
			wantErr: false,
		},
		{
			name: "Success test - Error from database",
			args: args{
				conn: db,
				data: RequestBody{
					Document: `654321`,
					Password: `some_pass`,
				},
			},
			wantFound: false,
			wantUser:  User{},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFound, gotUser, err := GetUserData(tt.args.conn, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFound != tt.wantFound {
				t.Errorf("GetUserData() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("GetUserData() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestGenerateJWT(t *testing.T) {
	type args struct {
		user User
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Success test",
			args: args{
				user: User{
					UserId: `123456`,
					Name:   `some_full_name`,
					Role:   `role`,
				},
			},
			want:    `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJpZCI6IjEyMzQ1NiIsIm5hbWUiOiJzb21lX2Z1bGxfbmFtZSIsInJvbGUiOiJyb2xlIn0.KWKl6ZT7HXpJpHmjISZ8yu0Yy-RxMeQkatCBXR35O30`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateJWT(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}
