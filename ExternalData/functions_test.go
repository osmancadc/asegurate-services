package main

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/go-sql-driver/mysql"
)

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
			name: "Success Test with error",
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
				t.Errorf("SuccessMessage() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestSetResponseHeaders(t *testing.T) {
	tests := []struct {
		name         string
		wantResponse events.APIGatewayProxyResponse
	}{
		{
			name: "Success Test",
			wantResponse: events.APIGatewayProxyResponse{
				Headers: map[string]string{
					"Content-Type":                 "application/json",
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Methods": "POST",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResponse := SetResponseHeaders(); !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("SetResponseHeaders() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestGetPersonData(t *testing.T) {
	os.Setenv(`BASE_URL`, `http://54.88.138.252:5000`)
	os.Setenv("AUTHORIZATION_TOKEN", "some_token")

	type args struct {
		data RequestGetData
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: "Success Test",
			args: args{
				data: RequestGetData{
					Document:       `123456`,
					ExpeditionDate: `01/01/2000`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{ "full_name":"some_full_name", "name":"some_name", "lastname":"some_lastname", "gender":"HOMBRE", "is_alive":true}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPersonData(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPersonData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want.StatusCode ||
				got.Body != tt.want.Body {
				t.Errorf("SuccessMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPersonName(t *testing.T) {
	os.Setenv(`BASE_URL`, `http://54.88.138.252:5000`)
	os.Setenv("AUTHORIZATION_TOKEN", "some_token")

	type args struct {
		data RequestGetName
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: "Success Test",
			args: args{
				data: RequestGetName{
					Document:     `123456`,
					DocumentType: `some_type`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"name":"some_name","lastname":"some_lastname"}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPersonName(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPersonName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want.StatusCode ||
				got.Body != tt.want.Body {
				t.Errorf("SuccessMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetProccedings(t *testing.T) {
	os.Setenv(`BASE_URL`, `http://54.88.138.252:5000`)
	os.Setenv("AUTHORIZATION_TOKEN", "some_token")
	type args struct {
		data RequestGetProccedings
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: `Success Test`,
			args: args{
				data: RequestGetProccedings{
					Document:     `123456`,
					DocumentType: `CC`,
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"formal_complaints":1,"recent_complain_year":1,"five_years_amount":1}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetProccedings(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProccedings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want.StatusCode ||
				got.Body != tt.want.Body {
				t.Errorf("GetProccedings() = %v, want %v", got, tt.want)
			}
		})
	}
}
