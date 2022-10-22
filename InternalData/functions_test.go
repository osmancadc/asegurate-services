package main

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
			name:    "Error Test",
			wantErr: true,
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

func TestErrorMessage(t *testing.T) {
	type args struct {
		functionError error
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: "Success Test",
			args: args{
				functionError: errors.New("some error"),
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some error"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := ErrorMessage(tt.args.functionError)
			if (err != nil) != tt.wantErr {
				t.Errorf("ErrorMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("SuccessMessage() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestSuccessMessage(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: "Success Test",
			args: args{
				message: "success message",
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"message":"success message"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := SuccessMessage(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("SuccessMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("SuccessMessage() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestSetResponseHeaders(t *testing.T) {
	tests := []struct {
		name         string
		wantResponse events.APIGatewayProxyResponse
	}{
		{
			name: `Success Test`,
			wantResponse: events.APIGatewayProxyResponse{
				Headers: map[string]string{
					"Content-Type":                 "application/json",
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "POST",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResponse := SetResponseHeaders(); !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("SetResponseHeaders() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestGetAuthorId(t *testing.T) {
	pool, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	columns := []string{`user_id`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`56`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`654321`).WillReturnRows(sqlmock.NewRows(columns))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`78910`)

	type args struct {
		conn     *gorm.DB
		document string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Success Test",
			args: args{
				conn:     db,
				document: `123456`,
			},
			want:    56,
			wantErr: false,
		},
		{
			name: "Error Test - User Not Found",
			args: args{
				conn:     db,
				document: `654321`,
			},
			want:    -1,
			wantErr: true,
		},
		{
			name: "Error Test - Database Error",
			args: args{
				conn:     db,
				document: `78910`,
			},
			want:    -1,
			wantErr: true,
		},
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

func TestGetPersonByDocument(t *testing.T) {
	pool, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(
			[]string{`name`, `gender`},
		).AddRow(`some_name`, `some_gender`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`1011112`).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{`name`, `lastname`},
			).AddRow(`some_name`, `some_lastname`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`654321`).
		WillReturnError(errors.New(`some_error`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`098765`).
		WillReturnRows(sqlmock.NewRows([]string{}))

	type args struct {
		conn *gorm.DB
		body GetByDocumentBody
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: `Success test - Get Name And Gender`,
			args: args{
				conn: db,
				body: GetByDocumentBody{
					Fields:   []string{`name`, `gender`},
					Document: `123456`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"name":"some_name","gender":"some_gender"}`,
			},
			wantErr: false,
		},
		{
			name: `Success test - Get Name And Lastname`,
			args: args{
				conn: db,
				body: GetByDocumentBody{
					Fields:   []string{`name`, `lastname`},
					Document: `1011112`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"name":"some_name","lastname":"some_lastname"}`,
			},
			wantErr: false,
		},
		{
			name: `Error test - Database Error`,
			args: args{
				conn: db,
				body: GetByDocumentBody{
					Document: `654321`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some_error"}`,
			},
			wantErr: false,
		},
		{
			name: `Error test - No Person Found`,
			args: args{
				conn: db,
				body: GetByDocumentBody{
					Document: `098765`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"no person found"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := GetPersonByDocument(tt.args.conn, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPersonByDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("GetPersonByDocument() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestGetScoreByDocument(t *testing.T) {
	pool, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	columns := []string{`col1`, `col2`, `col3`, `col4`}
	columns_error := []string{`col1`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(78.5, 5, 5, 45.0))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`654321`).
		WillReturnError(errors.New(`some database error`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`678901`).
		WillReturnRows(sqlmock.NewRows(columns_error).AddRow(78.5))
	type args struct {
		conn *gorm.DB
		body GetByDocumentBody
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: "Success Test",
			args: args{
				conn: db,
				body: GetByDocumentBody{
					Document: `123456`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{ "score": 78.500000, "positive_scores": 5, "negative_scores":5, "average_60_days":45.000000 }`,
			},
			wantErr: false,
		},
		{
			name: "Error Test - Database Error",
			args: args{
				conn: db,
				body: GetByDocumentBody{
					Document: `654321`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"internal score not found"}`,
			},
			wantErr: false,
		},
		{
			name: "Error Tests - Wrong Number Of Columns",
			args: args{
				conn: db,
				body: GetByDocumentBody{
					Document: `678901`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"internal score not found"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := GetScoreByDocument(tt.args.conn, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetScoreByDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("GetScoreByDocument() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestGetUserByPhone(t *testing.T) {
	pool, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	columns := []string{`document`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`3123456789`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`123456`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`3987654321`).
		WillReturnError(errors.New(`some error`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`302468024`).
		WillReturnRows(sqlmock.NewRows(columns))

	type args struct {
		conn *gorm.DB
		body GetByPhoneBody
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: `Success test`,
			args: args{
				conn: db,
				body: GetByPhoneBody{
					Phone: `3123456789`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"document":"123456"}`,
			},
			wantErr: false,
		},
		{
			name: `Error test - Database Error`,
			args: args{
				conn: db,
				body: GetByPhoneBody{
					Phone: `3987654321`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some error"}`,
			},
			wantErr: false,
		},
		{
			name: `Error test - User Not Found`,
			args: args{
				conn: db,
				body: GetByPhoneBody{
					Phone: `302468024`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"user not found"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := GetUserByPhone(tt.args.conn, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("GetUserByPhone() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestCheckUserByDocument(t *testing.T) {
	pool, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	columns := []string{`document`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`some_id`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`012345`).
		WillReturnRows(sqlmock.NewRows(columns))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`654321`).
		WillReturnError(errors.New(`some_error`))

	type args struct {
		conn *gorm.DB
		body GetByDocumentBody
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: "Success Test - User Exists",
			args: args{
				conn: db,
				body: GetByDocumentBody{
					Document: `123456`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"message":"user already exists"}`,
			},
		},
		{
			name: "Success Test - User Doesn't Exists",
			args: args{
				conn: db,
				body: GetByDocumentBody{
					Document: `012345`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"message":"user does not exists"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := CheckUserByDocument(tt.args.conn, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckUserByDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("CheckUserByDocument() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestInsertScore(t *testing.T) {
	pool, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	columns := []string{`user_id`}

	// First test mocks
	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(56))

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO  (.+)`).
		WithArgs(56, `678910`, 20, `test comment`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Second test mocks
	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`678901`).
		WillReturnRows(sqlmock.NewRows(columns))

	// Third test mocks
	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`098765`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(56))

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO  (.+)`).
		WithArgs(56, `678910`, 20, `test comment`).
		WillReturnError(errors.New(`some_error`))
	mock.ExpectRollback()

	type args struct {
		conn *gorm.DB
		body ScoreBody
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: "Success Test",
			args: args{
				conn: db,
				body: ScoreBody{
					Author:    `123456`,
					Objective: `678910`,
					Score:     20,
					Comments:  `test comment`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"message":"Score uploaded successfully"}`,
			},
		},
		{
			name: "Error Test - Author not found",
			args: args{
				conn: db,
				body: ScoreBody{
					Author:    `678901`,
					Objective: `678910`,
					Score:     20,
					Comments:  `test comment`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"author not found"}`,
			},
		},
		{
			name: "Error Test - Error From Database",
			args: args{
				conn: db,
				body: ScoreBody{
					Author:    `098765`,
					Objective: `678910`,
					Score:     20,
					Comments:  `test comment`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some_error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := InsertScore(tt.args.conn, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("InsertScore() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestInsertUser(t *testing.T) {
	pool, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//First test mock
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO  (.+)`).
		WithArgs(`123456`, `some@email.com`, `3123456789`, `some_password`, `some_role`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	//Second test mock
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO  (.+)`).
		WillReturnError(errors.New(`some_error`))
	mock.ExpectRollback()

	type args struct {
		conn *gorm.DB
		user User
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: `Success Test`,
			args: args{
				conn: db,
				user: User{
					Email:    `some@email.com`,
					Phone:    `3123456789`,
					Password: `some_password`,
					Document: `123456`,
					Role:     `some_role`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"message":"User inserted successfully"}`,
			},
			wantErr: false,
		},
		{
			name: `Error Test`,
			args: args{
				conn: db,
				user: User{},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some_error"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := InsertUser(tt.args.conn, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("InsertUser() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestInsertPerson(t *testing.T) {
	pool, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//First test mock
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO  (.+)`).
		WithArgs(`123456`, `some_name`, `some_lastname`, `some_gender`, 50, ``).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	//Second test mock
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO  (.+)`).
		WillReturnError(errors.New(`some_error`))
	mock.ExpectRollback()

	type args struct {
		conn   *gorm.DB
		person Person
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: "Success Test",
			args: args{
				conn: db,
				person: Person{
					Document: `123456`,
					Name:     `some_name`,
					Lastname: `some_lastname`,
					Gender:   `some_gender`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"message":"Person inserted successfully"}`,
			},
		},
		{
			name: "Error Test - Database Error",
			args: args{
				conn:   db,
				person: Person{},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some_error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := InsertPerson(tt.args.conn, tt.args.person)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertPerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("InsertPerson() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestUpdatePerson(t *testing.T) {
	pool, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//First test mock
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE (.+)`).
		WithArgs(`some_gender`, `123456`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	//Second test mock
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE (.+)`).
		WithArgs(`some_photo`, `123456`).
		WillReturnResult(sqlmock.NewResult(1, 0))
	mock.ExpectCommit()

	// //Third test mock
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE (.+)`).
		WithArgs(`some_photo`, `654321`).
		WillReturnError(errors.New(`some_error`))
	mock.ExpectRollback()

	type args struct {
		conn   *gorm.DB
		person Person
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: `Success Test - Update Gender`,
			args: args{
				conn: db,
				person: Person{
					Document: `123456`,
					Gender:   `some_gender`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"message":"Person data updated successfully"}`,
			},
		},
		{
			name: `Error Test - No Data Updated`,
			args: args{
				conn: db,
				person: Person{
					Photo:    `some_photo`,
					Document: `123456`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"no data was updated"}`,
			},
		},
		{
			name: `Error Test - Database Error`,
			args: args{
				conn: db,
				person: Person{
					Photo:    `some_photo`,
					Document: `654321`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some_error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := UpdatePerson(tt.args.conn, tt.args.person)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("UpdatePerson() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestGetAccountData(t *testing.T) {
	pool, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	columns := []string{`name`, `email`, `phone`, `photo`, `gender`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`some_name`, `some_email`, `some_phone`, `some_photo`, `some_gender`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnError(errors.New(`some_error`))

	type args struct {
		conn *gorm.DB
		body GetByDocumentBody
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: `Success Test`,
			args: args{
				conn: db,
				body: GetByDocumentBody{
					Document: `123456`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"name":"some_name","email":"some_email","phone":"some_phone","photo":"some_photo","gender":"some_gender"}`,
			},
			wantErr: false,
		},
		{
			name: `Error Test - Database Error`,
			args: args{
				conn: db,
				body: GetByDocumentBody{
					Document: `123456`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some_error"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := GetAccountData(tt.args.conn, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("GetAccountData() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestGetNameByPhone(t *testing.T) {
	pool, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	columns := []string{`name`, `lastname`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`31234444`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`some_name`, `some_lastname`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`3123456`).
		WillReturnRows(sqlmock.NewRows(columns))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`000`).
		WillReturnError(errors.New(`some_error`))

	type args struct {
		conn *gorm.DB
		body GetByPhoneBody
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: `Success Test`,
			args: args{
				conn: db,
				body: GetByPhoneBody{
					Phone: `31234444`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"name":"some_name","lastname":"some_lastname"}`,
			},
		},
		{
			name: `Error Test - No User Found`,
			args: args{
				conn: db,
				body: GetByPhoneBody{
					Phone: `3123456`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"no user found"}`,
			},
		},
		{
			name: `Error Test - Database Error`,
			args: args{
				conn: db,
				body: GetByPhoneBody{
					Phone: `000`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some_error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := GetNameByPhone(tt.args.conn, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNameByPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("GetNameByPhone() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestGetDocumentByPhone(t *testing.T) {
	pool, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	columns := []string{`document`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`31234444`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`123456`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`3123456`).
		WillReturnRows(sqlmock.NewRows(columns))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`000`).
		WillReturnError(errors.New(`some_error`))

	type args struct {
		conn *gorm.DB
		body GetByPhoneBody
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: `Success Test`,
			args: args{
				conn: db,
				body: GetByPhoneBody{
					Phone: `31234444`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"document":"123456"}`,
			},
		},
		{
			name: `Error Test - No Person Found`,
			args: args{
				conn: db,
				body: GetByPhoneBody{
					Phone: `3123456`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"no person found"}`,
			},
		},
		{
			name: `Error Test - Database Error`,
			args: args{
				conn: db,
				body: GetByPhoneBody{
					Phone: `000`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some_error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := GetDocumentByPhone(tt.args.conn, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDocumentByPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("GetDocumentByPhone() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}
