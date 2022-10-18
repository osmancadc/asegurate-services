package main

type RequestBody struct {
	Document       string `json:"document"`
	ExpeditionDate string `json:"expedition_date"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Role           string `json:"role"`
	Password       string `json:"password"`
}

type Person struct {
	Data PersonData `json:"data"`
}

type PersonData struct {
	Name     string `json:"name"`
	Lastname string `json:"lastName"`
	Gender   string `json:"gender"`
	IsAlive  bool   `json:"is_alive"`
}

type FindByDocumentBody struct {
	Document       string `json:"document"`
	ExpeditionDate string `json:"expedition_date,omitempty"`
}

type PersonBody struct {
	Document string `json:"document"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Gender   string `json:"gender"`
}

type UserBody struct {
	Document string `json:"document"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type DataBody struct {
	Document       string `json:"document"`
	ExpeditionDate string `json:"expedition_date"`
}

type InvokeBody struct {
	Action        string             `json:"action"`
	FindPerson    FindByDocumentBody `json:"get_by_document_body,omitempty"`
	Person        PersonBody         `json:"person_body,omitempty"`
	User          UserBody           `json:"user_body,omitempty"`
	GetPersonData DataBody           `json:"data_body,omitempty"`
}

type InvokePayload struct {
	Body string `json:"body"`
}

type InvokeResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type ResponseMesssage struct {
	Message string `json:"message"`
}
