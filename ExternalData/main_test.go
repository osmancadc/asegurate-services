package main

import (
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandlerExternalScoreData(t *testing.T) {
	os.Setenv(`BASE_URL`, `http://54.88.138.252:5000`)
	os.Setenv("AUTHORIZATION_TOKEN", "some_token")
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
			name: "Success test - Get Data",
			args: args{
				req: events.APIGatewayProxyRequest{
					Body: `{"action":"getPersonData","get_data":{"document":"123456","expedition_date":"01/01/2000"}}`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{ "full_name":"some_full_name", "name":"some_name", "lastname":"some_lastname", "gender":"HOMBRE", "is_alive":true}`,
			},
			wantErr: false,
		},
		{
			name: "Success test - Get Name",
			args: args{
				req: events.APIGatewayProxyRequest{
					Body: `{"action":"getPersonName","get_name":{"document":"123456","document_type":"CC"}}`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"name":"some_name","last_name":"some_lastname"}`,
			},
			wantErr: false,
		},
		{
			name: "Error test - Wrong Action",
			args: args{
				req: events.APIGatewayProxyRequest{
					Body: `{"action":"some_action"}`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       `{"message":"not a valid action"}`,
			},
			wantErr: false,
		},
		{
			name: "Error test - Wrong Request",
			args: args{
				req: events.APIGatewayProxyRequest{
					Body: ``,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"unexpected end of JSON input"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandlerExternalScoreData(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandlerExternalScoreData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want.StatusCode ||
				got.Body != tt.want.Body {
				t.Errorf("SuccessMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
