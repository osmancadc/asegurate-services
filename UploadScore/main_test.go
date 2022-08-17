package UploadScore

import (
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHanderUploadScore(t *testing.T) {
	type args struct {
		req events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HanderUploadScore(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HanderUploadScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HanderUploadScore() = %v, want %v", got, tt.want)
			}
		})
	}
}
