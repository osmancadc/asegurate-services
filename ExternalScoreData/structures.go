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
	Scope    string         `json:"scope"`
	DataBody RequestGetData `json:"get_data,omitempty"`
	NameBody RequestGetName `json:"get_name,omitempty"`
}
