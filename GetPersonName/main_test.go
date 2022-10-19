package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

func TestHandlerGetPersonName(t *testing.T) {
	OldGetClient := GetClient
	defer func() { GetClient = OldGetClient }()

	GetClient = func() lambdaiface.LambdaAPI {
		return &MockGetPersonNameMain{}
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
			name: `Success Tests - User Already Exists`,
			args: args{
				req: events.APIGatewayProxyRequest{
					PathParameters: map[string]string{`type`: `CC`, `value`: `123456`},
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{ "message": "success","name":"some_name some_lastname"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandlerGetPersonName(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandlerGetPersonName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want.StatusCode ||
				got.Body != tt.want.Body {
				t.Errorf("HandlerGetPersonName() = %v, want %v", got, tt.want)
			}
		})
	}
}
