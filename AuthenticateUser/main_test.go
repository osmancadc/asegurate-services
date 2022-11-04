package main

// import (
// 	"database/sql"
// 	"net/http"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/aws/aws-lambda-go/events"
// )

// func TestHandlerAuthenticateUser(t *testing.T) {
// 	OldConnectDatabase := ConnectDatabase
// 	defer func() { ConnectDatabase = OldConnectDatabase }()

// 	ConnectDatabase = func() (connection *sql.DB, err error) {
// 		db, mock, err := sqlmock.New()
// 		if err != nil {
// 			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 		}

// 		columns := []string{`document`, `name`, `role`}

// 		mock.ExpectQuery(`SELECT (.+) FROM (.+)`).
// 			WithArgs(`123456`, `some_pass`).
// 			WillReturnRows(sqlmock.NewRows(columns).AddRow(`123456`, `some_name`, `role`))

// 		return db, nil
// 	}

// 	type args struct {
// 		req events.APIGatewayProxyRequest
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    events.APIGatewayProxyResponse
// 		wantErr bool
// 	}{
// 		{
// 			name: `Success Test User Found`,
// 			args: args{
// 				req: events.APIGatewayProxyRequest{
// 					Body: `{ "document":"123456", "password":"some_pass" }`,
// 				},
// 			},
// 			want: events.APIGatewayProxyResponse{
// 				StatusCode: http.StatusOK,
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := HandlerAuthenticateUser(tt.args.req)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("HandlerAuthenticateUser(1) error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got.StatusCode != tt.want.StatusCode {
// 				t.Errorf("HandlerAuthenticateUser(2) = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
