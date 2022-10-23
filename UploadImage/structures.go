package main

type RequestBody struct {
	Image    string `json:"image"`
	Name     string `json:"name"`
	Document string `json:"document"`
}

type InvokePayload struct {
	Body string `json:"body"`
}

type InvokeBody struct {
	Action string `json:"action"`
	Person Person `json:"person_body,omitempty"`
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
