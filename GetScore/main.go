// package main

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/aws/aws-lambda-go/events"
// 	"github.com/aws/aws-lambda-go/lambda"
// )

// func HanderGetScore(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	var reqBody RequestBody

// 	response := events.APIGatewayProxyResponse{
// 		Headers: map[string]string{
// 			"Content-Type":                 "application/json",
// 			"Access-Control-Allow-Origin":  "*",
// 			"Access-Control-Allow-Methods": "POST",
// 		},
// 	}

// 	err := json.Unmarshal([]byte(req.Body), &reqBody)
// 	if err != nil {
// 		response.StatusCode = http.StatusBadRequest
// 		return response, err
// 	}

// 	conn := ConnectDatabase()
// 	defer conn.Close()

// 	if reqBody.Type == `CC` {
// 		return CalculateScoreDocument(reqBody, conn)
// 	} else {
// 		return CalculateScorePhone(reqBody, conn)
// 	}
// }

// func main() {
// 	lambda.Start(HanderGetScore)
// }

package main

import "fmt"

func main() {
	conn := ConnectDatabase()
	defer conn.Close()

	document := "1018500888"
	documentType := "CC"

	score := Score{}

	name, lastname, err := GetAssociatedName(document, documentType)
	if err != nil {
		fmt.Printf(`CalculateScore(1): %s`, err.Error())
		// return score, err
	}

	reputation := 50
	scoreAverage := 50
	var scoreArray []int

	results, err := conn.Query(`SELECT score FROM score WHERE objective = ?`, document)
	if err != nil {
		fmt.Printf(`CalculateScore(2): %s`, err.Error())
		// return score, err
	}

	for results.Next() {
		auxScore := 0
		err = results.Scan(&auxScore)
		if err != nil {
			fmt.Printf(`CalculateScore(3): %s`, err.Error())
			// return score, err
		}
		scoreArray = append(scoreArray, auxScore)
	}

	if len(scoreArray) > 0 {
		sum := 0
		for _, i := range scoreArray {
			sum += i
		}
		scoreAverage = int(sum / len(scoreArray))
	}

	fmt.Printf("%d -> %d = %s %s", reputation, scoreAverage, name, lastname)
	fmt.Println()

	fmt.Printf("-> %v \n", score)

	// return score, nil
}
