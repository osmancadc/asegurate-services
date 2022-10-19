package main

import (
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestHandlerInternalData(t *testing.T) {
	OldConnectDatabase := ConnectDatabase
	defer func() { ConnectDatabase = OldConnectDatabase }()

	ConnectDatabase = func() (db *gorm.DB, err error) {
		pool, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		db, err = gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{})
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
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
					Body: `{"action":"some"}`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       `{"message":"not a valid action"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandlerInternalData(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandlerInternalScoreData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Body != tt.want.Body ||
				got.StatusCode != tt.want.StatusCode {
				t.Errorf("HandlerInternalScoreData() = %v, want %v", got, tt.want)
			}
		})
	}
}
