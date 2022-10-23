package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandlerUploadScore(t *testing.T) {
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
			name: `Error Test - No Environment Variables`,
			args: args{
				req: events.APIGatewayProxyRequest{
					Body: `{"image":"","name":"some_name.txt","document":"123456"}`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"MissingRegion: could not find region configuration"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandlerUploadScore(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandlerUploadScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want.StatusCode ||
				got.Body != tt.want.Body {
				t.Errorf("HandlerUploadScore() = %v, want %v", got, tt.want)
			}
		})
	}
}
