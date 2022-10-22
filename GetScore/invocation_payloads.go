package main

import (
	"encoding/json"
	"time"
)

func GetDocumentInvokePayload(phone string) (payload []byte) {
	uploadBody, _ := json.Marshal(InvokeBody{
		Action: `getUserByPhone`,
		GetByPhone: GetByPhoneBody{
			Phone:  phone,
			Fields: []string{`document`},
		},
	})

	body := InvokePayload{
		Body: string(uploadBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func GetInternalScoreInvokePayload(document string) (payload []byte) {
	uploadBody, _ := json.Marshal(InvokeBody{
		Action: `getScore`,
		GetByDocument: GetByDocumentBody{
			Document: document,
		},
	})

	body := InvokePayload{
		Body: string(uploadBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func GetExternalProccedingsInvokePayload(document string) (payload []byte) {
	proccedingsBody, _ := json.Marshal(InvokeBody{
		Action: `getPersonProccedings`,
		GetExternalBody: GetExternalBody{
			Document: document,
			Type:     `CC`,
		},
	})

	body := InvokePayload{
		Body: string(proccedingsBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func GetStoredReputationInvokePayload(document string) (payload []byte) {
	uploadBody, _ := json.Marshal(InvokeBody{
		Action: `getPersonByDocument`,
		GetByDocument: GetByDocumentBody{
			Document: document,
			Fields:   []string{`reputation`, `last_update`, `name`, `lastname`, `gender`, `photo`},
		},
	})

	body := InvokePayload{
		Body: string(uploadBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func GetPredictScoreInvokePayload(internalScore InternalScore, externalProccedings ExternalProccedings) (payload []byte) {
	predictBody, _ := json.Marshal(PredictScoreBody{
		InternalScore:          internalScore.Score,
		InternalPositiveScores: internalScore.PositiveScores,
		InternalNegativeScores: internalScore.NegativeScores,
		InternalAverage60Days:  internalScore.Average60Days,
		FormalComplaints:       externalProccedings.FormalComplaints,
		FormalRecentYear:       externalProccedings.FormalRecentYear,
		Formal5YearsTotal:      externalProccedings.Formal5YearsTotal,
	})

	body := InvokePayload{
		Body: string(predictBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func UpdateSavedReputationInvokePayload(document string, reputation int) (payload []byte) {
	uploadBody, _ := json.Marshal(InvokeBody{
		Action: `updatePerson`,
		Person: Person{
			Document:   document,
			Reputation: reputation,
			LastUpdate: time.Now(),
		},
	})

	body := InvokePayload{
		Body: string(uploadBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func GetAssociatedNameInvokePayload(document string) (payload []byte) {
	getNameBody, _ := json.Marshal(InvokeBody{
		Action: `getPersonName`,
		GetExternalBody: GetExternalBody{
			Document: document,
			Type:     `CC`,
		},
	})

	body := InvokePayload{
		Body: string(getNameBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func SaveNewReputationInvokePayload(name, lastname, document string, reputation int) (payload []byte) {

	savePersonBody, _ := json.Marshal(InvokeBody{
		Action: `insertPerson`,
		Person: Person{
			Document:   document,
			Name:       name,
			Lastname:   lastname,
			Reputation: reputation,
		}})

	body := InvokePayload{
		Body: string(savePersonBody),
	}

	payload, _ = json.Marshal(body)

	return
}
