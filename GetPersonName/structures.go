package main

type Data struct {
	FullName  string `json:"fullName"`
	FirstName string `json:"firstName"`
	Lastname  string `json:"lastName"`
}

type Person struct {
	Data Data `json:"data"`
}
