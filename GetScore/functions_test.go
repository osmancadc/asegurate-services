package main

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
)

func TestValidatePhone(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{"document", "name", "lastname", "score", "reputation", "stars", "last_update"}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`3001234567`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`1018500888`, `Osman`, `Beltran Murcia`, `50`, `50`, `3`, "25-06-2022"))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`3007654321`).
		WillReturnRows(sqlmock.NewRows(columns))

	type args struct {
		conn  *sql.DB
		phone string
	}
	tests := []struct {
		name    string
		args    args
		want    Score
		want1   string
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				conn:  db,
				phone: `3001234567`,
			},
			want: Score{
				Name:       "Osman",
				Lastname:   "Beltran Murcia",
				Score:      50,
				Reputation: 50,
				Stars:      3,
				Updated:    `25-06-2022`,
			},
			want1:   `1018500888`,
			wantErr: false,
		},
		{
			name: "test not found",
			args: args{
				conn:  db,
				phone: `3007654321`,
			},
			want:    Score{},
			want1:   ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidatePhone(tt.args.conn, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidatePhone() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidatePhone() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetStoredScore(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{`name`, `lastname`, `score`, `reputation`, `stars`, `last_update`}

	mock.ExpectQuery(`SELECT (.+) FROM person (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`Osman`, `Beltran Murcia`, `50`, `50`, `3`, "25-06-2022"))

	mock.ExpectQuery(`SELECT (.+) FROM person (.+)`).
		WithArgs(`654321`).
		WillReturnRows(sqlmock.NewRows(columns))

	type args struct {
		conn     *sql.DB
		document string
	}
	tests := []struct {
		name    string
		args    args
		want    Score
		want1   bool
		wantErr bool
	}{
		{
			name: `Success test - User found`,
			args: args{
				conn:     db,
				document: `123456`,
			},
			want: Score{
				Name:       `Osman`,
				Lastname:   `Beltran Murcia`,
				Score:      50,
				Reputation: 50,
				Stars:      3,
				Updated:    `25-06-2022`,
			},
			want1:   true,
			wantErr: false,
		},
		{
			name: `Success test - No user found`,
			args: args{
				conn:     db,
				document: `654321`,
			},
			want:    Score{},
			want1:   false,
			wantErr: false,
		},
		{
			name: `Error test - No document sent`,
			args: args{
				conn: db,
			},
			want:    Score{},
			want1:   false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetStoredScore(tt.args.conn, tt.args.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStoredScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStoredScore() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetStoredScore() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDaysSinceLastUpdate(t *testing.T) {
	date := time.Now()

	type args struct {
		lastUpdate string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: `Success test - Correct date`,
			args: args{
				lastUpdate: date.Format(`2006-01-02 15:04:05`),
			},
			want:    0,
			wantErr: false,
		},
		{
			name: `Success test - Empty date`,
			args: args{
				lastUpdate: ``,
			},
			want:    1,
			wantErr: false,
		},
		{
			name: `Error test - Wrong date format`,
			args: args{
				lastUpdate: `01/01/2022-15HH:04M:05S`,
			},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DaysSinceLastUpdate(tt.args.lastUpdate)
			if (err != nil) != tt.wantErr {
				t.Errorf("DaysSinceLastUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DaysSinceLastUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAssociatedName(t *testing.T) {
	os.Setenv(`DATA_URL`, `https://asegurate2.free.beeceptor.com`)
	os.Setenv(`AUTHORIZATION_TOKEN`, `some-testing-token`)

	type args struct {
		document     string
		documentType string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: `Success test - User exists`,
			args: args{
				document:     `12345678`,
				documentType: `CC`,
			},
			want:    `OSMAN`,
			want1:   `BELTRAN MURCIA`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetAssociatedName(tt.args.document, tt.args.documentType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAssociatedName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAssociatedName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetAssociatedName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetResponseBody(t *testing.T) {
	type args struct {
		score    Score
		document string
		photo    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Success test - ",
			args: args{
				score: Score{
					Name:       "Osman",
					Lastname:   "Beltran",
					Score:      100,
					Reputation: 100,
					Stars:      5,
				},
				document: `1018500888`,
				photo:    `https://photo.jpg`,
			},
			want: fmt.Sprintf(`{
		"name": "%s",
		"document": "%s",
		"stars": %d,
		"reputation": %d,
		"score": %d,
		"certified": %t,
		"photo": "%s"
	}`, "Osman Beltran", "1018500888", 5, 100, 100, true, "https://photo.jpg"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetResponseBody(tt.args.score, tt.args.document, tt.args.photo); got != tt.want {
				t.Errorf("GetResponseBody() = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}

func TestSaveNewPerson(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare(`INSERT INTO person \((.+)\)`)
	mock.ExpectExec(`INSERT INTO person (.+)`).
		WithArgs(`123456`, `some_name`, `some_lastname`, 50, 3, 50).
		WillReturnResult(sqlmock.NewResult(1, 1))

	type args struct {
		conn     *sql.DB
		score    Score
		document string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success test - All the score data",
			args: args{
				conn: db,
				score: Score{
					Name:       `some_name`,
					Lastname:   `some_lastname`,
					Score:      50,
					Reputation: 50,
					Stars:      3,
				},
				document: `123456`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveNewPerson(tt.args.conn, tt.args.score, tt.args.document); (err != nil) != tt.wantErr {
				t.Errorf("SaveNewPerson() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalculateInternalScore(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{"score"}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`90`).AddRow(`10`).AddRow(`90`).AddRow(`50`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`654321`).
		WillReturnRows(sqlmock.NewRows(columns))

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
		{
			name: `Success test - 3 scores`,
			args: args{
				conn:     db,
				document: `123456`,
			},
			want:    60,
			wantErr: false,
		},
		{
			name: `Success test - No scores found`,
			args: args{
				conn:     db,
				document: `654321`,
			},
			want:    50,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateInternalScore(tt.args.conn, tt.args.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateInternalScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculateInternalScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPersonPhoto(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{"photo"}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`https://photo.png`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`3001234567`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`https://photo.png`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`3001234568`).
		WillReturnRows(sqlmock.NewRows(columns))

	type args struct {
		conn      *sql.DB
		dataValue string
		dataType  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: `Success test - using document`,
			args: args{
				conn:      db,
				dataValue: `123456`,
				dataType:  `CC`,
			},
			want:    `https://photo.png`,
			wantErr: false,
		},
		{
			name: `Success test - using cellphone`,
			args: args{
				conn:      db,
				dataValue: `3001234567`,
				dataType:  `PHONE`,
			},
			want:    `https://photo.png`,
			wantErr: false,
		},
		{
			name: `Success test - not photo found`,
			args: args{
				conn:      db,
				dataValue: `3001234568`,
				dataType:  `PHONE`,
			},
			want:    ``,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPersonPhoto(tt.args.conn, tt.args.dataValue, tt.args.dataType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPersonPhoto() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPersonPhoto() = %v, want %v", got, tt.want)
			}
		})
	}
}
