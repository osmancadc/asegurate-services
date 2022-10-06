package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

// TODO: Update Getdatabase mock function to increase percentage of coverage
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandlerGetScore(tt.args.req)
			t.Logf("Got: %v", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("HanderGetScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
