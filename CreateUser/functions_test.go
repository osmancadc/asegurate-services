package main

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	_ "github.com/go-sql-driver/mysql"
)

func TestGetClient(t *testing.T) {
	tests := []struct {
		name string
		want lambdaiface.LambdaAPI
	}{
		{
			name: `Single test`,
			want: &lambda.Lambda{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetClient()
		})
	}
}

func TestSetResponseHeaders(t *testing.T) {
	tests := []struct {
		name         string
		wantResponse events.APIGatewayProxyResponse
	}{
		{
			name: `Success Test`,
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

func TestCheckExistingPerson(t *testing.T) {
	type args struct {
		document string
		client   lambdaiface.LambdaAPI
	}
	tests := []struct {
		name            string
		args            args
		wantName        string
		wantExists      bool
		wantNeedsUpdate bool
		wantErr         bool
	}{
		{
			name: `Success Test - Status 200`,
			args: args{
				document: `123456`,
				client:   &MockCheckExistingPerson{},
			},
			wantName:        `some_name`,
			wantExists:      true,
			wantNeedsUpdate: false,
			wantErr:         false,
		},
		{
			name: `Success Test - Status Not 200`,
			args: args{
				document: ``,
				client:   &MockCheckExistingPerson{},
			},
			wantName:        ``,
			wantExists:      false,
			wantNeedsUpdate: false,
			wantErr:         false,
		},
		{
			name: `Error Test`,
			args: args{
				document: ``,
				client:   &MockCheckExistingPerson{},
			},
			wantName:        ``,
			wantExists:      false,
			wantNeedsUpdate: false,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotExists, gotNeedsUpdate, err := CheckExistingPerson(tt.args.document, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckExistingPerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("CheckExistingPerson() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotExists != tt.wantExists {
				t.Errorf("CheckExistingPerson() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
			if gotNeedsUpdate != tt.wantNeedsUpdate {
				t.Errorf("CheckExistingPerson() gotNeedsUpdate = %v, want %v", gotNeedsUpdate, tt.wantNeedsUpdate)
			}
		})
	}
}

func TestGetPersonData(t *testing.T) {

	type args struct {
		document       string
		expeditionDate string
		client         lambdaiface.LambdaAPI
	}
	tests := []struct {
		name           string
		args           args
		wantPersonData PersonData
		wantErr        bool
	}{
		{
			name: `Success Test`,
			args: args{
				document:       `123456`,
				expeditionDate: `01/01/2000`,
				client:         &MockGetPersonData{},
			},
			wantPersonData: PersonData{
				Name:     `some_name`,
				Lastname: `some_lastname`,
				Gender:   `some_gender`,
				IsAlive:  true,
			},
			wantErr: false,
		},
		{
			name: `Success Test - No data Returned`,
			args: args{
				document:       `123456`,
				expeditionDate: `01/01/2000`,
				client:         &MockGetPersonData{},
			},
			wantPersonData: PersonData{},
			wantErr:        false,
		},
		{
			name: `error Test - Invocation Error`,
			args: args{
				document:       `123456`,
				expeditionDate: `01/01/2000`,
				client:         &MockGetPersonData{},
			},
			wantPersonData: PersonData{},
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPersonData, err := GetPersonData(tt.args.document, tt.args.expeditionDate, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPersonData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPersonData, tt.wantPersonData) {
				t.Errorf("GetPersonData() = %v, want %v", gotPersonData, tt.wantPersonData)
			}
		})
	}
}

func TestSavePerson(t *testing.T) {
	type args struct {
		personData PersonData
		document   string
		client     lambdaiface.LambdaAPI
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: `Success Test`,
			args: args{
				personData: PersonData{
					Name:     `some_name`,
					Lastname: `some_lastname`,
					Gender:   `some_gender`,
				},
				document: `123456`,
				client:   &MockSavePerson{},
			},
			wantErr: false,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				personData: PersonData{
					Gender: `HOMBRE`,
				},
				document: `123456`,
				client:   &MockSavePerson{},
			},
			wantErr: true,
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				personData: PersonData{},
				document:   `123456`,
				client:     &MockSavePerson{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SavePerson(tt.args.personData, tt.args.document, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("SavePerson() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateGender(t *testing.T) {
	type args struct {
		gender   string
		document string
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
				gender:   `some_gender`,
				document: `123456`,
				client:   &MockUpdateGender{},
			},
			wantErr: false,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				gender: `HOMBRE`,
				client: &MockUpdateGender{},
			},
			wantErr: true,
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				client: &MockUpdateGender{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateGender(tt.args.gender, tt.args.document, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("UpdateGender() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePersonData(t *testing.T) {
	type args struct {
		personData PersonData
		auxErr     error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: `Success Test`,
			args: args{
				personData: PersonData{
					Name:    `some_name`,
					IsAlive: true,
				},
				auxErr: nil,
			},
			wantErr: false,
		},
		{
			name: `Error Test - Empty Name`,
			args: args{
				personData: PersonData{
					Name:    ``,
					IsAlive: true,
				},
				auxErr: nil,
			},
			wantErr: true,
		},
		{
			name: `Error Test - Dead Person`,
			args: args{
				personData: PersonData{
					Name:    `some_name`,
					IsAlive: false,
				},
				auxErr: nil,
			},
			wantErr: true,
		},
		{
			name: `Error Test - AuxError Not Null`,
			args: args{
				personData: PersonData{},
				auxErr:     errors.New(`some_error`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatePersonData(tt.args.personData, tt.args.auxErr); (err != nil) != tt.wantErr {
				t.Errorf("ValidatePersonData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsertPerson(t *testing.T) {
	type args struct {
		document       string
		expeditionDate string
		client         lambdaiface.LambdaAPI
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantErr  bool
	}{
		{
			name: `Success Test - Person Already Exists`,
			args: args{

				client: &MockInsertPerson{},
			},
			wantName: `some_name`,
			wantErr:  false,
		},
		{
			name: `Error Test - Person Is Dead`,
			args: args{
				client: &MockInsertPerson{},
			},
			wantName: ``,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, err := InsertPerson(tt.args.document, tt.args.expeditionDate, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertPerson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("InsertPerson() = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}

func TestCheckExistingUser(t *testing.T) {
	type args struct {
		document string
		client   lambdaiface.LambdaAPI
	}
	tests := []struct {
		name       string
		args       args
		wantExists bool
		wantErr    bool
	}{
		{
			name: `Success Test - User Already Exists`,
			args: args{
				document: `123456`,
				client:   &MockCheckExistingUser{},
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			name: `Success Test - User Doesn't Exists`,
			args: args{
				document: `123456`,
				client:   &MockCheckExistingUser{},
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				client: &MockCheckExistingUser{},
			},
			wantExists: false,
			wantErr:    true,
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				client: &MockCheckExistingUser{},
			},
			wantExists: false,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExists, err := CheckExistingUser(tt.args.document, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckExistingUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExists != tt.wantExists {
				t.Errorf("CheckExistingUser() = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}

func TestSaveUser(t *testing.T) {
	type args struct {
		user   UserBody
		client lambdaiface.LambdaAPI
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: `Success Test`,
			args: args{
				user: UserBody{
					Document: `123456`,
					Email:    `some@email.com`,
					Phone:    `3123456789`,
					Password: `some_password`,
					Role:     `some_role`,
				},
				client: &MockSaveUser{},
			},
			wantErr: false,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				user:   UserBody{},
				client: &MockSaveUser{},
			},
			wantErr: true,
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				user:   UserBody{},
				client: &MockSaveUser{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveUser(tt.args.user, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("SaveUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsertUser(t *testing.T) {
	type args struct {
		user   UserBody
		client lambdaiface.LambdaAPI
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: `Success Test`,
			args: args{
				user: UserBody{
					Document: `123456`,
					Email:    `some@email.com`,
					Phone:    `3123456789`,
					Password: `some_password`,
					Role:     `some_role`,
				},
				client: &MockInsertUser{},
			},
			wantErr: false,
		},
		{
			name: `Error Test - User Already Exists`,
			args: args{
				user: UserBody{
					Document: `123456`,
					Email:    `some@email.com`,
					Phone:    `3123456789`,
					Password: `some_password`,
					Role:     `some_role`,
				},
				client: &MockInsertUser{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertUser(tt.args.user, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("InsertUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
