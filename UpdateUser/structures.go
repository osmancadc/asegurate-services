package main

type RequestBody struct {
	Image    string `json:"image"`
	Name     string `json:"name"`
	Document string `json:"document"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

type InvokePayload struct {
	Body string `json:"body"`
}

type InvokeBody struct {
	Action string `json:"action"`
	Person Person `json:"person_body,omitempty"`
	User   User   `json:"user_body,omitempty"`
}

type InvokeResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

type Person struct {
	Document string `json:"document,omitempty"`
	Photo    string `json:"photo,omitempty"`
}

type User struct {
	Document string `json:"document,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
}
