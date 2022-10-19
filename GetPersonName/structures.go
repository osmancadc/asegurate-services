package main

type Data struct {
	FullName  string `json:"fullName"`
	FirstName string `json:"firstName"`
	Lastname  string `json:"lastName"`
}

type Person struct {
	Data Data `json:"data"`
}

type InvokePayload struct {
	Body string `json:"body"`
}

type InvokeBody struct {
	Action          string            `json:"action"`
	GetByDocument   GetByDocumentBody `json:"get_by_document_body,omitempty"`
	GetByPhone      GetByPhoneBody    `json:"get_by_phone_body,omitempty"`
	GetExternalBody GetExternalBody   `json:"name_body,omitempty"`
}

type GetExternalBody struct {
	Document     string `json:"document"`
	DocumentType string `json:"document_type"`
}

type GetByDocumentBody struct {
	Document string   `json:"document"`
	Fields   []string `json:"fields"`
}

type GetByPhoneBody struct {
	Phone string `json:"phone"`
}

type InvokeResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type ResponseBody struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
}

type MessageBody struct {
	Message string `json:"message"`
}
