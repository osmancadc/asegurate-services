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
	type args struct {
		score Score
	}
	tests := []struct {
		name         string
		args         args
		wantResponse events.APIGatewayProxyResponse
		wantErr      bool
	}{
		{
			name: `Success Test`,
			args: args{
				score: Score{
					Document:   `123456`,
					Name:       `some_name`,
					Gender:     `some_gender`,
					Score:      0,
					Reputation: 0,
					Stars:      0,
				},
			},
			wantResponse: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"document":"123456","name":"some_name","gender":"some_gender","score":0,"reputation":0,"stars":0,"certified":false,"photo":""}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := SuccessMessage(tt.args.score)
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

func TestGetAssociatedDocument(t *testing.T) {
	type args struct {
		phone  string
		client lambdaiface.LambdaAPI
	}
	tests := []struct {
		name         string
		args         args
		wantDocument string
	}{
		{
			name: `Success Test`,
			args: args{
				phone:  `3123456`,
				client: &MockGetAssociatedDocument{},
			},
			wantDocument: `123456`,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				phone:  `3123456`,
				client: &MockGetAssociatedDocument{},
			},
			wantDocument: ``,
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				phone:  `3123456`,
				client: &MockGetAssociatedDocument{},
			},
			wantDocument: ``,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDocument := GetAssociatedDocument(tt.args.phone, tt.args.client); gotDocument != tt.wantDocument {
				t.Errorf("GetAssociatedDocument() = %v, want %v", gotDocument, tt.wantDocument)
			}
		})
	}
}

func TestGetInternalScore(t *testing.T) {
	type args struct {
		document string
		client   lambdaiface.LambdaAPI
	}
	tests := []struct {
		name         string
		args         args
		wantScore    InternalScore
		wantIsStored bool
		wantErr      bool
	}{
		{
			name: `Success Test`,
			args: args{
				document: `123456`,
				client:   &MockGetInternalScore{},
			},
			wantScore: InternalScore{
				Score:          0,
				PositiveScores: 0,
				NegativeScores: 0,
				Average60Days:  0,
			},
			wantIsStored: true,
			wantErr:      false,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				document: ``,
				client:   &MockGetInternalScore{},
			},
			wantScore:    InternalScore{},
			wantIsStored: false,
			wantErr:      true,
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				document: ``,
				client:   &MockGetInternalScore{},
			},
			wantScore:    InternalScore{},
			wantIsStored: false,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotScore, gotIsStored, err := GetInternalScore(tt.args.document, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInternalScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotScore, tt.wantScore) {
				t.Errorf("GetInternalScore() gotScore = %v, want %v", gotScore, tt.wantScore)
			}
			if gotIsStored != tt.wantIsStored {
				t.Errorf("GetInternalScore() gotIsStored = %v, want %v", gotIsStored, tt.wantIsStored)
			}
		})
	}
}

func TestGetExternalProccedings(t *testing.T) {
	type args struct {
		document string
		client   lambdaiface.LambdaAPI
	}
	tests := []struct {
		name            string
		args            args
		wantProccedings ExternalProccedings
		wantErr         bool
	}{
		{
			name: `Success Test`,
			args: args{
				document: `123456`,
				client:   &MockGetExternalProceedings{},
			},
			wantProccedings: ExternalProccedings{
				FormalComplaints:  0,
				FormalRecentYear:  0,
				Formal5YearsTotal: 0,
			},
			wantErr: false,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				document: ``,
				client:   &MockGetExternalProceedings{},
			},
			wantProccedings: ExternalProccedings{
				FormalComplaints:  0,
				FormalRecentYear:  0,
				Formal5YearsTotal: 0,
			},
			wantErr: true,
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				document: ``,
				client:   &MockGetExternalProceedings{},
			},
			wantProccedings: ExternalProccedings{
				FormalComplaints:  0,
				FormalRecentYear:  0,
				Formal5YearsTotal: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProccedings, err := GetExternalProccedings(tt.args.document, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExternalProccedings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotProccedings, tt.wantProccedings) {
				t.Errorf("GetExternalProccedings() = %v, want %v", gotProccedings, tt.wantProccedings)
			}
		})
	}
}

func TestGetStoredScore(t *testing.T) {
	type args struct {
		document string
		client   lambdaiface.LambdaAPI
	}
	tests := []struct {
		name            string
		args            args
		wantStoredScore Score
	}{
		{
			name: `Success Test`,
			args: args{
				document: `123456`,
				client:   &MockGetStoredScore{},
			},
			wantStoredScore: Score{
				Name:       `some_name some_lastname`,
				Reputation: 0,
			},
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				document: ``,
				client:   &MockGetStoredScore{},
			},
			wantStoredScore: Score{
				Name:       ``,
				Reputation: 0,
			},
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				document: ``,
				client:   &MockGetStoredScore{},
			},
			wantStoredScore: Score{
				Name:       ``,
				Reputation: 0,
			},
		},
		{
			name: `Error Test - Empty Response`,
			args: args{
				document: ``,
				client:   &MockGetStoredScore{},
			},
			wantStoredScore: Score{
				Name:       ``,
				Reputation: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStoredScore, _ := GetStoredScore(tt.args.document, tt.args.client)
			if !reflect.DeepEqual(gotStoredScore, tt.wantStoredScore) {
				t.Errorf("GetStoredScore() gotStoredScore = %v, want %v", gotStoredScore, tt.wantStoredScore)
			}

		})
	}
}

func TestPredictPersonScore(t *testing.T) {
	type args struct {
		internalScore       InternalScore
		externalProccedings ExternalProccedings
		client              lambdaiface.LambdaAPI
	}
	tests := []struct {
		name                   string
		args                   args
		wantPredictionResponse PredictionResponse
		wantErr                bool
	}{
		{
			name: `Success Test`,
			args: args{
				internalScore:       InternalScore{},
				externalProccedings: ExternalProccedings{},
				client:              &MockPredictScore{},
			},
			wantPredictionResponse: PredictionResponse{},
			wantErr:                false,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				internalScore:       InternalScore{},
				externalProccedings: ExternalProccedings{},
				client:              &MockPredictScore{},
			},
			wantPredictionResponse: PredictionResponse{},
			wantErr:                true,
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				internalScore:       InternalScore{},
				externalProccedings: ExternalProccedings{},
				client:              &MockPredictScore{},
			},
			wantPredictionResponse: PredictionResponse{},
			wantErr:                false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPredictionResponse, err := PredictPersonScore(tt.args.internalScore, tt.args.externalProccedings, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("PredictPersonScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPredictionResponse, tt.wantPredictionResponse) {
				t.Errorf("PredictPersonScore() = %v, want %v", gotPredictionResponse, tt.wantPredictionResponse)
			}
		})
	}
}

func TestCalculateScore(t *testing.T) {
	type args struct {
		document string
		client   lambdaiface.LambdaAPI
	}
	tests := []struct {
		name         string
		args         args
		wantIsStored bool
		wantScore    Score
		wantErr      bool
	}{
		{
			name: `Success Test`,
			args: args{
				document: `123456`,
				client:   &MockCalculateScore{},
			},
			wantIsStored: true,
			wantScore: Score{
				Document:  `123456`,
				Name:      `some_name some_lastname`,
				Certified: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsStored, gotScore, err := CalculateScore(tt.args.document, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsStored != tt.wantIsStored {
				t.Errorf("CalculateScore() gotIsStored = %v, want %v", gotIsStored, tt.wantIsStored)
			}
			if !reflect.DeepEqual(gotScore, tt.wantScore) {
				t.Errorf("CalculateScore() gotScore = %v, want %v", gotScore, tt.wantScore)
			}
		})
	}
}

func TestUpdateSavedReputation(t *testing.T) {
	type args struct {
		document   string
		reputation int
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
				client: &MockUpdate{},
			},
			wantErr: false,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				client: &MockUpdate{},
			},
			wantErr: true,
		},
		{
			name: `Error Test - Status 500`,
			args: args{
				client: &MockUpdate{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateSavedReputation(tt.args.document, tt.args.reputation, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("UpdateSavedReputation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAssociatedName(t *testing.T) {
	type args struct {
		document string
		client   lambdaiface.LambdaAPI
	}
	tests := []struct {
		name         string
		args         args
		wantName     string
		wantLastname string
		wantErr      bool
	}{
		{
			name: `Success Tests`,
			args: args{
				client: &MockAssociatedName{},
			},
			wantName:     `some_name`,
			wantLastname: `some_lastname`,
			wantErr:      false,
		},
		{
			name: `Error Tests - Invocation Error`,
			args: args{
				client: &MockAssociatedName{},
			},
			wantName:     ``,
			wantLastname: ``,
			wantErr:      true,
		},
		{
			name: `Error Tests - Status 500`,
			args: args{
				client: &MockAssociatedName{},
			},
			wantName:     ``,
			wantLastname: ``,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotLastname, err := GetAssociatedName(tt.args.document, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAssociatedName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("GetAssociatedName() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotLastname != tt.wantLastname {
				t.Errorf("GetAssociatedName() gotLastname = %v, want %v", gotLastname, tt.wantLastname)
			}
		})
	}
}

func TestSaveNewReputation(t *testing.T) {
	type args struct {
		document   string
		reputation int
		client     lambdaiface.LambdaAPI
	}
	tests := []struct {
		name         string
		args         args
		wantName     string
		wantLastname string
		wantErr      bool
	}{
		{
			name: `Success Test`,
			args: args{
				client: &MockCreateScore{},
			},
			wantName:     `some_name`,
			wantLastname: `some_lastname`,
			wantErr:      false,
		},
		{
			name: `Error Test - Invocation Error`,
			args: args{
				client: &MockCreateScore{},
			},
			wantName:     `some_name`,
			wantLastname: `some_lastname`,
			wantErr:      true,
		},
		{
			name: `Success Test - Status 500`,
			args: args{
				client: &MockCreateScore{},
			},
			wantName:     `some_name`,
			wantLastname: `some_lastname`,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotLastname, err := SaveNewReputation(tt.args.document, tt.args.reputation, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveNewReputation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("SaveNewReputation() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotLastname != tt.wantLastname {
				t.Errorf("SaveNewReputation() gotLastname = %v, want %v", gotLastname, tt.wantLastname)
			}
		})
	}
}
