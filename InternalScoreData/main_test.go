package main

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
)

func TestHandlerInternalScoreData(t *testing.T) {
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
			WillReturnRows(sqlmock.NewRows(columns).AddRow(`56`))

		mock.ExpectPrepare(`INSERT INTO score \((.+)\)`)

		mock.ExpectExec(`INSERT INTO score \((.+)\)`).
			WithArgs(56, `654321`, 99, `some comment`).
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
			name: "Insert Test",
			args: args{
				req: events.APIGatewayProxyRequest{
					Body: `{ "action": "insert", "insert_data": { "author": "123456", "objective": "654321" ,"score": 99,"comments": "some comment" } }`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"message":"Score uploaded successfully"}`,
			},
		},
		{
			name: "Test error no a valid action",
			args: args{
				req: events.APIGatewayProxyRequest{
					Body: `{ "action": "some_weird_action" }`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       `{"message":"not a valid action"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandlerInternalScoreData(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandlerInternalScoreData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want.StatusCode ||
				got.Body != tt.want.Body {
				t.Errorf("HandlerInternalScoreData() = %v, want %v", got, tt.want)
			}
		})
	}
}
