package main

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandlerGetUserData(t *testing.T) {

	OldConnectDatabase := ConnectDatabase
	defer func() { ConnectDatabase = OldConnectDatabase }()

	ConnectDatabase = func() (connection *sql.DB, err error) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		columns := []string{`name`, `email`, `phone`, `photo`, `gender`}

		mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
			WithArgs(`123456`).
			WillReturnRows(sqlmock.NewRows(columns).AddRow(`some_full_name`, `some@email.com`, `300123456`, `http://photo.png`, `male`))

		return db, nil
	}

	pathParameter := make(map[string]string)
	pathParameter[`document`] = `123456`

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
			name: `Success Test HandlerGetUserData`,
			args: args{
				req: events.APIGatewayProxyRequest{
					PathParameters: pathParameter,
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
			got, err := HandlerGetUserData(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandlerGetUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want.StatusCode {
				t.Errorf("HandlerGetUserData() = %v, want %v", got, tt.want)
			}
		})
	}
}
