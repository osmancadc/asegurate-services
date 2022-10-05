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
			name:    "Success Test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ConnectDatabase()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConnectDatabase() error = %v, wantErr %v", err, tt.wantErr)
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

	columns := []string{`name`, `email`, `phone`, `photo`, `gender`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`some_full_name`, `some@email.com`, `300123456`, `http://photo.png`, `male`))

	type args struct {
		document string
		conn     *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    User
		wantErr bool
	}{
		{
			name: `Success Test`,
			args: args{
				document: `123456`,
				conn:     db,
			},
			want: User{
				Name:     `some_full_name`,
				Document: `123456`,
				Email:    `some@email.com`,
				Phone:    `300123456`,
				Photo:    `http://photo.png`,
				Gender:   `male`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserData(tt.args.document, tt.args.conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserData() = %v, want %v", got, tt.want)
			}
		})
	}
}
