package main

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

type MockFindByIdClientMain struct {
	lambdaiface.LambdaAPI
}

func (mlc *MockFindByIdClientMain) Invoke(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	payload, _ := json.Marshal(InvokeResponse{
		StatusCode: 200,
		Body:       `{"document":"654321"}`,
	})

	return &lambda.InvokeOutput{
		Payload: payload,
	}, nil
}

func TestHandlerUploadScore(t *testing.T) {

	OldGetClient := GetClient
	defer func() { GetClient = OldGetClient }()

	GetClient = func() lambdaiface.LambdaAPI {
		return &MockFindByIdClientMain{}
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
			wantErr: false,
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
			wantErr: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			_, err := HandlerUploadScore(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandlerUploadScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
