package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandlerGetScore(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var reqBody RequestBody

	err := json.Unmarshal([]byte(req.Body), &reqBody)
	if err != nil {
		response := SetResponseHeaders()
		response.StatusCode = http.StatusBadRequest
		return response, nil
	}

	client := GetClient()

	document := reqBody.Value

	if reqBody.Type != `CC` {
		document = GetAssociatedDocument(reqBody.Value, client)
	}

	if document == `` {
		return ErrorMessage(errors.New(`no user found`))
	}

	isStored, score, err := CalculateScore(document, client)
	if err != nil {
		return ErrorMessage(err)
	}

	if isStored {
		UpdateSavedReputation(document, score.Reputation, client)
	} else {
		name, lastname, err := SaveNewReputation(document, score.Reputation, client)
		if err != nil {
			return ErrorMessage(err)
		}
		score.Name = fmt.Sprintf(`%s %s`, name, lastname)
		score.Gender = `male`
	}

	score.Stars = int(math.Round(float64(score.Reputation+score.Score) / float64(40.0)))

	return SuccessMessage(score)

}

func main() {
	lambda.Start(HandlerGetScore)
}
