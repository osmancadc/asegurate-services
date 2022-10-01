package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHanderGetScore(t *testing.T) {

	type args struct {
		req events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Main test",
			args: args{
				events.APIGatewayProxyRequest{
					Body: `{ "value": "1018500888", "type": "CC" } `,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HanderGetScore(tt.args.req)
			t.Logf("Got: %v", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("HanderGetScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
