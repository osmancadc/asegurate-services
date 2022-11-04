package main

type RequestBody struct {
	Document string `json:"document"`
	Password string `json:"password"`
}

type User struct {
	UserId string
	Name   string
	Role   string
}

type InvokePayload struct {
	Body string `json:"body"`
}

type InvokeBody struct {
	Action        string            `json:"action"`
	GetByDocument GetByDocumentBody `json:"get_by_document_body,omitempty"`
}

type GetByDocumentBody struct {
	Document string   `json:"document"`
	Fields   []string `json:"fields"`
}

type InvokeResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type ResponseBody struct {
	Password string `json:"password"`
}

type MessageBody struct {
	Message string `json:"message"`
}
