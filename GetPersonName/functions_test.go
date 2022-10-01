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

func TestGetFromDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{`name`}
	columns_error := []string{`some`, `some2`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`some_name some_lastname`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+) INNER JOIN (.+)`).
		WithArgs(`3001234567`).
		WillReturnRows(sqlmock.NewRows(columns_error).AddRow(`some result`, `some other result`))

	type args struct {
		conn      *sql.DB
		dataType  string
		dataValue string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		want1   string
		wantErr bool
	}{
		{
			name: "Success Test Using document",
			args: args{
				conn:      db,
				dataType:  `CC`,
				dataValue: `123456`,
			},
			want:    true,
			want1:   `some_name some_lastname`,
			wantErr: false,
		},
		{
			name: "Error Test Incorrect number of columns in response",
			args: args{
				conn:      db,
				dataType:  `PHONE`,
				dataValue: `3001234567`,
			},
			want:    false,
			want1:   ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetFromDatabase(tt.args.conn, tt.args.dataType, tt.args.dataValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFromDatabase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFromDatabase() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetFromDatabase() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetFromProvider(t *testing.T) {
	os.Setenv(`DATA_URL`, `https://asegurate3.free.beeceptor.com`)
	os.Setenv(`AUTHORIZATION_TOKEN`, `some-testing-token`)
	type args struct {
		dataType  string
		dataValue string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		want1   string
		wantErr bool
	}{
		{
			name: "Success Test Provider Response",
			args: args{
				dataType:  `CC`,
				dataValue: `123456`,
			},
			want:    true,
			want1:   `OSMAN BELTRAN MURCIA`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetFromProvider(tt.args.dataType, tt.args.dataValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFromProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFromProvider() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetFromProvider() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
