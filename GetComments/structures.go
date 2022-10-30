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

type MessageBody struct {
	Message string `json:"message"`
}

type Comments struct {
	Comments []Comment `json:"comments"`
}

type Comment struct {
	Author  string `json:"author"`
	Photo   string `json:"photo"`
	Comment string `json:"comment"`
	Score   int    `json:"score"`
}
