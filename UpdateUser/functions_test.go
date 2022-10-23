package main

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

func TestGetClient(t *testing.T) {
	tests := []struct {
		name string
		want lambdaiface.LambdaAPI
	}{
		{
			name: `Success Test`,
			want: &lambda.Lambda{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetClient()
		})
	}
}

func TestErrorMessage(t *testing.T) {
	type args struct {
		functionError error
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: "Success Test",
			args: args{
				functionError: errors.New(`some error`),
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       `{"message":"some error"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := ErrorMessage(tt.args.functionError)
			if (err != nil) != tt.wantErr {
				t.Errorf("ErrorMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("ErrorMessage() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestSuccessMessage(t *testing.T) {
	tests := []struct {
		name         string
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: `Success Test`,
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"message":"image updated successfully"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := SuccessMessage()
			if (err != nil) != tt.wantErr {
				t.Errorf("SuccessMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResponse.StatusCode != tt.wantResponse.StatusCode ||
				gotResponse.Body != tt.wantResponse.Body {
				t.Errorf("SuccessMessage() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestGenerateName(t *testing.T) {
	type args struct {
		document  string
		date      string
		extension string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: `Success Tests`,
			args: args{
				document:  `123456`,
				date:      `some_date`,
				extension: `ext`,
			},
			want: `c98e40c9d24f70e77516f731d862044cdde5ef1ba45ffac6e921148ef2c687ec8a868f2f091f70e6f7c42c06b01b128b65ca50e07a88cca185424c7b55c209af.ext`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateName(tt.args.document, tt.args.date, tt.args.extension); got != tt.want {
				t.Errorf("GenerateName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveTemporalFile(t *testing.T) {
	type args struct {
		data     []byte
		name     string
		document string
	}
	tests := []struct {
		name                 string
		args                 args
		wantTemporalFileName string
		wantErr              bool
	}{
		{
			name: `Success Test`,
			args: args{
				data:     []byte{},
				name:     `some_name.ext`,
				document: `123456`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTemporalFileName, err := SaveTemporalFile(tt.args.data, tt.args.name, tt.args.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveTemporalFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			os.Remove(gotTemporalFileName)
		})
	}
}

func TestGetTemporalFile(t *testing.T) {
	type args struct {
		temporalName string
	}
	tests := []struct {
		name         string
		args         args
		wantFileName string
		wantErr      bool
	}{
		{
			name: `Error Test - Non Existing File`,
			args: args{
				temporalName: `some_name`,
			},
			wantFileName: ``,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotFileName, err := GetTemporalFile(tt.args.temporalName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTemporalFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotFileName != tt.wantFileName {
				t.Errorf("GetTemporalFile() gotFileName = %v, want %v", gotFileName, tt.wantFileName)
			}
		})
	}
}

func TestUploadToS3(t *testing.T) {
	type args struct {
		file     io.Reader
		fileName string
	}
	tests := []struct {
		name         string
		args         args
		wantLocation string
		wantErr      bool
	}{
		{
			name: `Error Test - No Environment Variables`,
			args: args{
				file:     &io.LimitedReader{},
				fileName: ``,
			},
			wantLocation: ``,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLocation, err := UploadToS3(tt.args.file, tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadToS3() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLocation != tt.wantLocation {
				t.Errorf("UploadToS3() = %v, want %v", gotLocation, tt.wantLocation)
			}
		})
	}
}

func TestUploadImage(t *testing.T) {
	type args struct {
		data     []byte
		name     string
		document string
	}
	tests := []struct {
		name         string
		args         args
		wantLocation string
		wantErr      bool
	}{
		{
			name: `Error Test - No Environment Variables`,
			args: args{
				data:     []byte{},
				name:     `some_name.ext`,
				document: `123456`,
			},
			wantLocation: ``,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLocation, err := UploadImage(tt.args.data, tt.args.name, tt.args.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLocation != tt.wantLocation {
				t.Errorf("UploadImage() = %v, want %v", gotLocation, tt.wantLocation)
			}
		})
	}
}

func TestUpdateDatabase(t *testing.T) {
	type args struct {
		route    string
		document string
		email    string
		phone    string
		client   lambdaiface.LambdaAPI
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: `Success Test`,
			args: args{
				client: &MockUpdateDatabase{},
			},
			wantErr: false,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				client: &MockUpdateDatabase{},
			},
			wantErr: true,
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				client: &MockUpdateDatabase{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateDatabase(tt.args.route, tt.args.document, tt.args.email, tt.args.phone, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("UpdateDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
