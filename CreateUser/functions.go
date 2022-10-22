package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	str "strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	invokeLambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"

	_ "github.com/go-sql-driver/mysql"
)

func GetClient() lambdaiface.LambdaAPI {
	region := os.Getenv(`REGION`)
	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	return invokeLambda.New(sess, &aws.Config{Region: aws.String(region)})
}

func SetResponseHeaders() (response events.APIGatewayProxyResponse) {
	response = events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		},
	}
	return
}

func ErrorMessage(functionError error) (response events.APIGatewayProxyResponse, err error) {
	response = SetResponseHeaders()

	response.StatusCode = http.StatusInternalServerError
	response.Body = fmt.Sprintf(`{"message":"%s"}`, functionError.Error())

	return
}

func GetInvokePayload(document, action string) (payload []byte) {

	findPersonBody, _ := json.Marshal(InvokeBody{Action: action, FindPerson: FindByDocumentBody{Document: document}})

	body := InvokePayload{
		Body: string(findPersonBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func ExternalDataInvokePayload(document, expeditionDate string) (payload []byte) {
	findPersonBody, _ := json.Marshal(InvokeBody{
		Action: `getPersonData`,
		GetPersonData: DataBody{
			Document:       document,
			ExpeditionDate: expeditionDate,
		}})

	body := InvokePayload{
		Body: string(findPersonBody),
	}

	payload, _ = json.Marshal(body)

	return
}

/*
Takes a `PersonData` struct and a document number as parameters, and returns a payload `[]byte` that can be used to invoke the lambda function
*/
func SavePersonInvokePayload(personData PersonData, document string) (payload []byte) {
	gender := `female`

	if personData.Gender == `HOMBRE` {
		gender = `male`
	}

	savePersonBody, _ := json.Marshal(InvokeBody{
		Action: `insertPerson`,
		Person: PersonBody{
			Document: document,
			Name:     personData.Name,
			Lastname: personData.Lastname,
			Gender:   gender,
		}})

	body := InvokePayload{
		Body: string(savePersonBody),
	}

	payload, _ = json.Marshal(body)

	return
}

/*
Takes a gender and a document as parameters, and returns a payload `[]byte` that can be used to invoke the lambda function
*/
func UpdateGenderInvokePayload(gender, document string) (payload []byte) {

	if gender == `HOMBRE` {
		gender = `male`
	} else {
		gender = `female`
	}

	updatePersonBody, _ := json.Marshal(InvokeBody{
		Action: `updatePerson`,
		Person: PersonBody{
			Document: document,
			Gender:   gender,
		}})

	body := InvokePayload{
		Body: string(updatePersonBody),
	}

	payload, _ = json.Marshal(body)

	return
}

/*
Takes a `UserBody` struct as parameter, and returns a payload `[]byte` that can be used to invoke the lambda function
*/
func SaveUserInvokePayload(user UserBody) (payload []byte) {

	saveUserBody, _ := json.Marshal(InvokeBody{
		Action: `insertUser`,
		User:   user,
	})

	body := InvokePayload{
		Body: string(saveUserBody),
	}

	payload, _ = json.Marshal(body)

	return
}

/*
Takes a document number and a lambda client, and returns the name of the person, a boolean
indicating if the person exists, a boolean indicating if the person needs to be updated, and an error
*/
func CheckExistingPerson(document string, client lambdaiface.LambdaAPI) (name string, exists bool, needsUpdate bool, err error) {
	payload := GetInvokePayload(document, `getPersonByDocument`)

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`CheckExistingPerson(1): %s`, err.Error())
		return
	}

	response := InvokeResponse{}
	json.Unmarshal(result.Payload, &response)

	if response.StatusCode != 200 {
		fmt.Printf(`CheckExistingPerson(2): %s`, response.Body)
		return
	}

	bodyString := str.Replace(string(response.Body), `\`, ``, -1)
	personData := PersonData{}
	json.Unmarshal([]byte(bodyString), &personData)

	if personData.Name != `not_found` {
		exists = true
		name = personData.Name
		if personData.Gender == `` {
			needsUpdate = true
		}
	}

	return
}

func GetPersonData(document, expeditionDate string, client lambdaiface.LambdaAPI) (personData PersonData, err error) {
	payload := ExternalDataInvokePayload(document, expeditionDate)

	some := string(payload)
	fmt.Println(some)

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("ExternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`GetPersonData(1): %s`, err.Error())
		return
	}

	response := InvokeResponse{}
	json.Unmarshal(result.Payload, &response)

	if response.StatusCode != 200 {
		fmt.Printf(`GetPersonData(2): %s`, response.Body)
		return
	}

	bodyString := str.Replace(string(response.Body), `\`, ``, -1)
	json.Unmarshal([]byte(bodyString), &personData)

	return
}

func SavePerson(personData PersonData, document string, client lambdaiface.LambdaAPI) (err error) {
	payload := SavePersonInvokePayload(personData, document)

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`SavePerson(1): %s`, err.Error())
		return
	}

	response := InvokeResponse{}
	json.Unmarshal(result.Payload, &response)

	if response.StatusCode != 200 {
		responseMessage := ResponseMesssage{}
		bodyString := str.Replace(string(response.Body), `\`, ``, -1)
		json.Unmarshal([]byte(bodyString), &responseMessage)

		fmt.Printf(`SavePerson(2): %s`, responseMessage.Message)
		err = errors.New(responseMessage.Message)
		return
	}

	return
}

func UpdateGender(gender, document string, client lambdaiface.LambdaAPI) (err error) {
	payload := UpdateGenderInvokePayload(gender, document)

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`UpdateGender(1): %s`, err.Error())
		return
	}

	response := InvokeResponse{}
	json.Unmarshal(result.Payload, &response)

	if response.StatusCode != 200 {
		responseMessage := ResponseMesssage{}
		bodyString := str.Replace(string(response.Body), `\`, ``, -1)
		json.Unmarshal([]byte(bodyString), &responseMessage)

		fmt.Printf(`UpdateGender(2): %s`, responseMessage.Message)
		err = errors.New(responseMessage.Message)
		return
	}

	return
}

func ValidatePersonData(personData PersonData, auxErr error) (err error) {
	if auxErr != nil {
		err = auxErr
		fmt.Printf(`ValidatePersonData(1): %s`, err.Error())
		return
	}

	if personData.Name == `` {
		fmt.Printf(`ValidatePersonData(2): Person not found`)
		err = errors.New(`no se encontro a la persona, intenta nuevamente`)
		return
	}

	if !personData.IsAlive {
		fmt.Printf(`ValidatePersonData(3): Person is dead`)
		err = errors.New(`el documento es de una persona fallecida`)
		return
	}

	return
}

func InsertPerson(document, expeditionDate string, client lambdaiface.LambdaAPI) (name string, err error) {
	name, exists, needsUpdate, err := CheckExistingPerson(document, client)
	if err != nil {
		fmt.Printf(`InsertPerson(1): %s`, err.Error())
		return
	}
	if exists && !needsUpdate {
		fmt.Println(`InsertPerson(2): Person already exists`)
		return
	}

	personData, err := GetPersonData(document, expeditionDate, client)
	err = ValidatePersonData(personData, err)
	if err != nil {
		fmt.Printf(`InsertPerson(3): %s`, err.Error())
		return
	}

	if needsUpdate {
		err = UpdateGender(personData.Gender, document, client)
		if err != nil {
			fmt.Printf(`InsertPerson(4): %s`, err.Error())
			return
		}
	} else {
		err = SavePerson(personData, document, client)
		if err != nil {
			fmt.Printf(`InsertPerson(4): %s`, err.Error())
			return
		}

	}

	return personData.Name, nil
}

func CheckExistingUser(document string, client lambdaiface.LambdaAPI) (exists bool, err error) {
	payload := GetInvokePayload(document, `getPersonByDocument`)
	response := InvokeResponse{}
	responseMesssage := ResponseMesssage{}

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`CheckExistingUser(1): %s`, err.Error())
		return
	}

	json.Unmarshal(result.Payload, &response)
	bodyString := str.Replace(string(response.Body), `\`, ``, -1)
	json.Unmarshal([]byte(bodyString), &responseMesssage)

	if response.StatusCode != 200 {
		err = errors.New(responseMesssage.Message)
		fmt.Printf(`CheckExistingUser(2): %s`, response.Body)
		return
	}

	exists = (responseMesssage.Message == `user already exists`)
	return
}

func SaveUser(user UserBody, client lambdaiface.LambdaAPI) (err error) {
	payload := SaveUserInvokePayload(user)

	result, err := client.Invoke(&invokeLambda.InvokeInput{FunctionName: aws.String("InternalData"), Payload: payload})
	if err != nil {
		fmt.Printf(`SaveUser(1): %s`, err.Error())
		return
	}

	response := InvokeResponse{}
	json.Unmarshal(result.Payload, &response)

	if response.StatusCode != 200 {
		responseMessage := ResponseMesssage{}
		bodyString := str.Replace(string(response.Body), `\`, ``, -1)
		json.Unmarshal([]byte(bodyString), &responseMessage)

		fmt.Printf(`SaveUser(2): %s`, responseMessage.Message)
		err = errors.New(responseMessage.Message)
		return
	}

	return
}

func InsertUser(user UserBody, client lambdaiface.LambdaAPI) (err error) {

	exists, err := CheckExistingUser(user.Document, client)
	if err != nil {
		fmt.Printf(`InsertUser(1): %s`, err.Error())
		return
	}

	if exists {
		fmt.Println(`InsertUser(2): User already exists`)
		err = errors.New(`user already exists`)
		return
	}

	err = SaveUser(user, client)

	return
}
