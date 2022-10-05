package main

import (
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
)

func TestHandlerCreateUser(t *testing.T) {
	os.Setenv(`DATA_URL`, `http://54.88.138.252:5000`)
	os.Setenv(`AUTHORIZATION_TOKEN`, `some-testing-token`)

	OldConnectDatabase := ConnectDatabase
	defer func() { ConnectDatabase = OldConnectDatabase }()

	ConnectDatabase = func() (connection *sql.DB, err error) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		columns := []string{`name`, `email`, `phone`, `photo`}

		mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
			WithArgs(`123456`).
			WillReturnRows(sqlmock.NewRows(columns))

		mock.ExpectPrepare(`INSERT INTO person \((.+)\)`)
		mock.ExpectExec(`INSERT INTO person \((.+)\)`).
			WithArgs(`123456`, `some_name`, `some_lastname`).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
			WithArgs(`123456`).
			WillReturnRows(sqlmock.NewRows(columns))

		mock.ExpectPrepare(`INSERT INTO user \((.+)\)`)
		mock.ExpectExec(`INSERT INTO user \((.+)\)`).
			WithArgs(`some@email.com`, `3001234567`, `some_password`, `123456`, `role`).
			WillReturnResult(sqlmock.NewResult(1, 1))

		return db, nil
	}

	type args struct {
		req events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: `Success Test`,
			args: args{
				req: events.APIGatewayProxyRequest{
					Body: `{
						"document": "123456",
						"expiration_date":"06/02/1998",
						"email": "some@email.com",
						"phone": "3001234567",
						"role": "role",
						"password": "some_password"
					}`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandlerCreateUser(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HanderCreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want.StatusCode {
				t.Errorf("HanderCreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
