package main

// import (
// 	"os"
// 	"testing"

// 	"github.com/aws/aws-lambda-go/events"
// )

// // TODO: Update Getdatabase mock function to increase percentage of coverage
// func TestHandlerGetPersonName(t *testing.T) {
// 	os.Setenv(`DATA_URL`, `http://54.88.138.252:5000`)
// 	os.Setenv(`AUTHORIZATION_TOKEN`, `some-testing-token`)
// 	type args struct {
// 		req events.APIGatewayProxyRequest
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name:    "Handler testing with empty request",
// 			args:    args{events.APIGatewayProxyRequest{}},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, err := HandlerGetPersonName(tt.args.req)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("HandlerGetPersonName() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 		})
// 	}
// }
