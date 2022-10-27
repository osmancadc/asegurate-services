package main

type RequestBody struct {
	Document string `json:"document"`
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
	Email string `json:"email"`
}

type MessageBody struct {
	Message string `json:"message"`
}
