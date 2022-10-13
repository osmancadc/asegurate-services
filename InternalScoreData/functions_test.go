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

func TestUploadInternalScore(t *testing.T) {
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
		body InsertBody
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
				body: InsertBody{
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
				body: InsertBody{
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
			gotResponse, err := UploadInternalScore(tt.args.conn, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadInternalScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("UploadInternalScore() = %v, want %v", gotResponse, tt.wantResponse)
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
		body GetBody
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
				body: GetBody{
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
				body: GetBody{
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
				body: GetBody{
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
		body UpdateBody
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
				body: UpdateBody{
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
