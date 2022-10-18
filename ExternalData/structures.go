package main

type RequestBody struct {
	Action          string                `json:"action"`
	DataBody        RequestGetData        `json:"data_body,omitempty"`
	NameBody        RequestGetName        `json:"name_body,omitempty"`
	ProccedingsBody RequestGetProccedings `json:"proccedings_body,omitempty"`
}

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

type RequestGetProccedings struct {
	Document     string `json:"document"`
	DocumentType string `json:"document_type"`
}

type ProccedingsResponse struct {
	Data ProccedingsData `json:"data"`
}

type ProccedingsData struct {
	Subject     string            `json:"consultedSubject"`
	Proccedings []ProccedingsList `json:"list"`
	Record      Record            `json:"pagination"`
}

type ProccedingsList struct {
	Department string `json:"departamento"`
	Date       string `json:"fechaUltimaActuacion"`
	IsPrivate  bool   `json:"esPrivado"`
}

type Record struct {
	Total int `json:"records"`
}
