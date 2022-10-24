package main

type User struct {
	Document string `json:"document"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Photo    string `json:"photo"`
	Gender   string `json:"gender"`
	Role     string `json:"role"`
}
type InvokePayload struct {
	Body string `json:"body"`
}

type InvokeBody struct {
	Action      string            `json:"action"`
	GetUserData GetByDocumentBody `json:"get_by_document_body,omitempty"`
}

type GetByDocumentBody struct {
	Document string `json:"document"`
}

type InvokeResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}
