package main

type RequestBody struct {
	Document string `json:"document"`
	Type     string `json:"type"`
}

type Score struct {
	Name       string
	Lastname   string
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
