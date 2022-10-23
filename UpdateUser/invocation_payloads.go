package main

import (
	"encoding/json"
)

func UpdatePhotoInvokePayload(document, photo string) (payload []byte) {
	uploadBody, _ := json.Marshal(InvokeBody{
		Action: `updatePerson`,
		Person: Person{
			Document: document,
			Photo:    photo,
		},
	})

	body := InvokePayload{
		Body: string(uploadBody),
	}

	payload, _ = json.Marshal(body)

	return
}

func UpdateUserInvokePayload(document, email, phone string) (payload []byte) {
	uploadBody, _ := json.Marshal(InvokeBody{
		Action: `updateUser`,
		User: User{
			Document: document,
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
