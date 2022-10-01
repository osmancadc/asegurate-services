package main

import (
	"database/sql"
	"os"
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

func TestUploadScorePhone(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{`document`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`3001234567`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`123456`))

	mock.ExpectPrepare(`INSERT INTO score \((.+)\)`)
	mock.ExpectExec(`INSERT INTO score \((.+)\)`).
		WithArgs(50, `123456`, 50, `No comments`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`3001234568`).
		WillReturnRows(sqlmock.NewRows(columns))

	type args struct {
		conn     *sql.DB
		author   int
		score    int
		phone    string
		comments string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success test - Valid cellphone",
			args: args{
				conn:     db,
				author:   50,
				score:    50,
				phone:    `3001234567`,
				comments: `No comments`,
			},
			wantErr: false,
		},
		{
			name: "Success test - Invalid cellphone",
			args: args{
				conn:  db,
				phone: `3001234568`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UploadScorePhone(tt.args.conn, tt.args.author, tt.args.score, tt.args.phone, tt.args.comments); (err != nil) != tt.wantErr {
				t.Errorf("UploadScorePhone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAuthorId(t *testing.T) {
	type args struct {
		conn     *sql.DB
		document string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAuthorId(tt.args.conn, tt.args.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAuthorId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAuthorId() = %v, want %v", got, tt.want)
			}
		})
	}
}
