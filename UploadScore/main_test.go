package main

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
)

func TestHanderUploadScore(t *testing.T) {
	OldConnectDatabase := ConnectDatabase
	defer func() { ConnectDatabase = OldConnectDatabase }()

	ConnectDatabase = func() (connection *sql.DB, err error) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		columns := []string{`user_id`}

		// First test mocks

		mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
			WithArgs(`123456`).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(`50`))

		mock.ExpectPrepare(`INSERT INTO score \((.+)\)`)
		mock.ExpectExec(`INSERT INTO score \((.+)\)`).
			WithArgs(50, `123456`, 50, `No comments`).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Second test mocks

		mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
			WithArgs(`123456`).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(`50`))

		mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
			WithArgs(`3001234567`).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(`123456`))
		mock.ExpectPrepare(`INSERT INTO score \((.+)\)`)

		mock.ExpectExec(`INSERT INTO score \((.+)\)`).
			WithArgs(50, `123456`, 50, `No comments`).
			WillReturnResult(sqlmock.NewResult(1, 1))

		return db, nil
	}

	type args struct {
		req events.APIGatewayProxyRequest
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func() (connection *sql.DB, err error)
		wantErr  bool
	}{
		{
			name: "Handler test with CC in the request",
			args: args{
				req: events.APIGatewayProxyRequest{
					Body: `{
						"author": "123456",
						"type": "CC",
						"value": "123456",
						"score": 50,
						"comments": "No comments"
					}`,
				},
			},
			mockFunc: ConnectDatabase,
			wantErr:  false,
		},
		{
			name: "Handler test with PHONE in the request",
			args: args{
				req: events.APIGatewayProxyRequest{
					Body: `{
						"author": "123456",
						"type": "PHONE",
						"value": "3001234567",
						"score": 50,
						"comments": "No comments"
					}`,
				},
			},
			mockFunc: ConnectDatabase,
			wantErr:  false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			_, err := HanderUploadScore(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HanderUploadScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
