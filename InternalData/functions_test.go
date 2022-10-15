package main

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
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
			name: "Testing with message",
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
			name: "Testing with message",
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

func TestGetAuthorId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{`col1`}
	columns_error := []string{`col1`, `col2`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`56`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`654321`).WillReturnRows(sqlmock.NewRows(columns))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`78910`)

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`09876`).WillReturnRows(sqlmock.NewRows(columns_error).AddRow(`56`, `65`))

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
			name: "Tests success",
			args: args{
				conn:     db,
				document: `123456`,
			},
			want:    56,
			wantErr: false,
		},
		{
			name: "Test no user found",
			args: args{
				conn:     db,
				document: `654321`,
			},
			want:    -1,
			wantErr: true,
		},
		{
			name: "Test error in database",
			args: args{
				conn:     db,
				document: `78910`,
			},
			want:    -1,
			wantErr: true,
		},
		{
			name: "Test error in return of database",
			args: args{
				conn:     db,
				document: `09876`,
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

func TestInsertInternalScore(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{`col1`}
	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`56`))

	mock.ExpectPrepare(`INSERT INTO score \((.+)\)`)

	mock.ExpectExec(`INSERT INTO score \((.+)\)`).
		WithArgs(56, `678910`, 20, `test comment`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`678901`).
		WillReturnRows(sqlmock.NewRows(columns))

	type args struct {
		conn *sql.DB
		body InsertScoreBody
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: "Tests success",
			args: args{
				conn: db,
				body: InsertScoreBody{
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
			name: "Tests error author not found",
			args: args{
				conn: db,
				body: InsertScoreBody{
					Author:    `678901`,
					Objective: `678910`,
					Score:     20,
					Comments:  `test comment`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"no user found"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := InsertInternalScore(tt.args.conn, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertInternalScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("InsertInternalScore() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestGetInternalScoreSummary(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

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
		conn *sql.DB
		body GetScoreBody
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: "Test succesfull",
			args: args{
				conn: db,
				body: GetScoreBody{
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
			name: "Test error database result",
			args: args{
				conn: db,
				body: GetScoreBody{
					Document: `654321`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some database error"}`,
			},
			wantErr: false,
		},
		{
			name: "Test error database number of columns",
			args: args{
				conn: db,
				body: GetScoreBody{
					Document: `678901`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"sql: expected 1 destination arguments in Scan, not 4"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := GetInternalScoreSummary(tt.args.conn, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInternalScoreSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("InsertInternalScore() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestUpdateInternalScore(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare(`UPDATE person SET (.+)`)

	mock.ExpectExec(`UPDATE person (.+)`).
		WithArgs(56, 20, `123456`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	type args struct {
		conn *sql.DB
		body UpdateScoreBody
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: "Test successfull",
			args: args{
				conn: db,
				body: UpdateScoreBody{
					Document:   `123456`,
					Score:      56,
					Reputation: 20,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"message":"User score updated successfully"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := UpdateInternalScore(tt.args.conn, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateInternalScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("InsertInternalScore() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestGetUserByPhone(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{`col1`}
	columns_error := []string{`col1`, `col2`}
	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`3123456789`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`123456`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`3987654321`).
		WillReturnError(errors.New(`some error`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456789`).
		WillReturnRows(sqlmock.NewRows(columns_error).AddRow(`123456`, `something`))

	type args struct {
		conn *sql.DB
		body GetUserByPhoneBody
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
				body: GetUserByPhoneBody{
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
			name: `Error test - database error`,
			args: args{
				conn: db,
				body: GetUserByPhoneBody{
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
			name: `Error test - wrong database return`,
			args: args{
				conn: db,
				body: GetUserByPhoneBody{
					Phone: `123456789`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"sql: expected 2 destination arguments in Scan, not 1"}`,
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
				t.Errorf("InsertInternalScore() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestGetPersonByDocument(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{`col1`, `col2`}
	columns_error := []string{`col1`}

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`123456`).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(`some_name`, `some_gender`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`654321`).
		WillReturnError(errors.New(`some_error`))

	mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
		WithArgs(`654321`).
		WillReturnRows(sqlmock.NewRows(columns_error).AddRow(`some_name`))

	type args struct {
		conn *sql.DB
		body GetByDocumentBody
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
				body: GetByDocumentBody{
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
			name: `Error test - Error from DB`,
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
			name: `Error test - Wrong amount of cols`,
			args: args{
				conn: db,
				body: GetByDocumentBody{
					Document: `654321`,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"sql: expected 1 destination arguments in Scan, not 2"}`,
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

func TestGetUserByDocument(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{`col1`}

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
		conn *sql.DB
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
			name: "Success Test - User Doesn's Exists",
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
				Body:       `{"message":"some_error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := GetUserByDocument(tt.args.conn, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("GetUserByDocument() = %v, want %v", gotResponse, tt.wantResponse)
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

	// columns := []string{`col1`}

	mock.ExpectPrepare(`INSERT INTO user (.+)`)

	mock.ExpectExec(`INSERT INTO user (.+)`).
		WithArgs(`some@email.com`, `3123456789`, `some_password`, `123456`, `some_role`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectPrepare(`INSERT INTO user (.+)`).WillReturnError(errors.New(`some_error`))

	mock.ExpectPrepare(`INSERT INTO user (.+)`)

	mock.ExpectExec(`INSERT INTO user (.+)`).
		WithArgs(``, ``, ``, ``, ``).
		WillReturnError(errors.New(`some_error`))

	type args struct {
		conn *sql.DB
		body InsertUserBody
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
				body: InsertUserBody{
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
			name: `Error Test - Error Preparing Query`,
			args: args{
				conn: db,
				body: InsertUserBody{},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some_error"}`,
			},
			wantErr: false,
		},
		{
			name: `Error Test - Error Executing Query`,
			args: args{
				conn: db,
				body: InsertUserBody{},
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
			gotResponse, err := InsertUser(tt.args.conn, tt.args.body)
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

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare(`INSERT INTO person (.+)`)

	mock.ExpectExec(`INSERT INTO person (.+)`).
		WithArgs(`123456`, `some_name`, `some_lastname`, `some_gender`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectPrepare(`INSERT INTO person (.+)`).WillReturnError(errors.New(`some_error`))

	mock.ExpectPrepare(`INSERT INTO person (.+)`)
	mock.ExpectExec(`INSERT INTO person (.+)`).
		WillReturnError(errors.New(`some_error`))

	type args struct {
		conn *sql.DB
		body InsertPersonBody
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
				body: InsertPersonBody{
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
			name: "Error Test - Preparing Statement",
			args: args{
				conn: db,
				body: InsertPersonBody{},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some_error"}`,
			},
		},
		{
			name: "Error Test - Executing Query",
			args: args{
				conn: db,
				body: InsertPersonBody{},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some_error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := InsertPerson(tt.args.conn, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertPerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("InsertUser() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}
