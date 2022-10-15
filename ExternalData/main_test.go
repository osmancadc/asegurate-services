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
					Body: `{"scope":"data","get_data":{"document":"1018500888","expedition_date":"01/01/2000"}}`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{ "full_name":"some_full_name", "first_name":"some_name", "last_name":"some_lastname", "gender":"HOMBRE", "isAlive":true}`,
			},
			wantErr: false,
		},
		{
			name: "Success test - Get Name",
			args: args{
				req: events.APIGatewayProxyRequest{
					Body: `{"scope":"name","get_name":{"document":"1018500888","document_type":"CC"}}`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"name":"some_name","last_name":"some_lastname"}`,
			},
			wantErr: false,
		},
		{
			name: "Error test - Bad Scope",
			args: args{
				req: events.APIGatewayProxyRequest{
					Body: `{"scope":"some_scope"}`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       `{"message":"not a valid scope"}`,
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
