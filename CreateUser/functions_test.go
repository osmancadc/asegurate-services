package main

import (
	"database/sql"
	"errors"
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

func TestCheckExistingUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{`document`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`123456`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`654321`).
		WillReturnRows(sqlmock.NewRows(columns))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`98765`).
		WillReturnError(errors.New("some error"))

	type args struct {
		conn     *sql.DB
		document string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: `Success Test - User already exists`,
			args: args{
				conn:     db,
				document: `123456`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: `Success Test - User doesn't exists`,
			args: args{
				conn:     db,
				document: `654321`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: `Error Test Database Error`,
			args: args{
				conn:     db,
				document: `98765`,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckExistingUser(tt.args.conn, tt.args.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckExistingUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckExistingUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPersonData(t *testing.T) {
	os.Setenv(`DATA_URL`, `http://54.88.138.252:5000`)
	os.Setenv(`AUTHORIZATION_TOKEN`, `some-testing-token`)
	type args struct {
		document       string
		expirationDate string
	}
	tests := []struct {
		name    string
		args    args
		want    PersonData
		wantErr bool
	}{
		{
			name: `Success Test Provider`,
			args: args{
				document:       `123456`,
				expirationDate: `23/03/2022`,
			},
			want: PersonData{
				Name:     `some_name`,
				Lastname: `some_lastname`,
				IsAlive:  true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPersonData(tt.args.document, tt.args.expirationDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPersonData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPersonData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertPerson(t *testing.T) {
	os.Setenv(`DATA_URL`, `http://54.88.138.252:5000`)
	os.Setenv(`AUTHORIZATION_TOKEN`, `some-testing-token`)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{`name`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns))

	mock.ExpectPrepare(`INSERT INTO person \((.+)\)`)
	mock.ExpectExec(`INSERT INTO person \((.+)\)`).
		WithArgs(`123456`, `some_name`, `some_lastname`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock second test

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`some_name`))

	type args struct {
		conn           *sql.DB
		document       string
		expirationDate string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: `Success Test`,
			args: args{
				conn:           db,
				document:       `123456`,
				expirationDate: `23/08/2022`,
			},
			want:    `some_name`,
			wantErr: false,
		},
		{
			name: `Success Test User Already Exists`,
			args: args{
				conn:           db,
				document:       `123456`,
				expirationDate: `23/08/2022`,
			},
			want:    `some_name`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InsertPerson(tt.args.conn, tt.args.document, tt.args.expirationDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertPerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("InsertPerson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{`name`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns))

	mock.ExpectPrepare(`INSERT INTO user \((.+)\)`)
	mock.ExpectExec(`INSERT INTO user \((.+)\)`).
		WithArgs(`some@email.com`, `3001234567`, `some_password`, `123456`, `role`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	type args struct {
		conn     *sql.DB
		email    string
		phone    string
		password string
		document string
		role     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: `Success Test`,
			args: args{
				conn:     db,
				email:    `some@email.com`,
				phone:    `3001234567`,
				password: `some_password`,
				document: `123456`,
				role:     `role`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertUser(tt.args.conn, tt.args.email, tt.args.phone, tt.args.password, tt.args.document, tt.args.role); (err != nil) != tt.wantErr {
				t.Errorf("InsertUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
