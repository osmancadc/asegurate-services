package main

type RequestBody struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Score struct {
	Name       string
	Lastname   string
	Gender     string
	Score      int
	Reputation int
	Stars      int
	Updated    string
}

type Data struct {
	FullName  string `json:"fullName"`
	FirstName string `json:"firstName"`
	Lastname  string `json:"lastName"`
}

type Person struct {
	Data Data `json:"data"`
}
