package main

type PersonData struct {
	FullName string `json:"fullName"`
	Name     string `json:"firstName"`
	Lastname string `json:"lastName"`
	IsAlive  bool   `json:"isAlive"`
	Gender   string `json:"gender"`
}

type Person struct {
	Data PersonData `json:"data"`
}

type RequestGetData struct {
	Document       string `json:"document"`
	ExpeditionDate string `json:"expedition_date"`
}

type RequestGetName struct {
	Document     string `json:"document"`
	DocumentType string `json:"document_type"`
}

type RequestBody struct {
	Action   string         `json:"action"`
	DataBody RequestGetData `json:"data_body,omitempty"`
	NameBody RequestGetName `json:"name_body,omitempty"`
}
