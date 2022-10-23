package main

import (
	"encoding/json"
)

func UpdateSavedReputationInvokePayload(document, photo, email, phone string) (payload []byte) {
	uploadBody, _ := json.Marshal(InvokeBody{
		Action: `updatePerson`,
		Person: Person{
			Document: document,
			Photo:    photo,
			Email:    email,
			Phone:    phone,
		},
	})

	body := InvokePayload{
		Body: string(uploadBody),
	}

	payload, _ = json.Marshal(body)

	return
}
