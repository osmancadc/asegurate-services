package main

type RequestBody struct {
	Document       string `json:"document"`
	ExpirationDate string `json:"expiration_date"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Role           string `json:"role"`
	Password       string `json:"password"`
}

type Person struct {
	Data PersonData `json:"data"`
}

type PersonData struct {
	Name     string `json:"firstName"`
	Lastname string `json:"lastName"`
	IsAlive  bool   `json:"isAlive"`
}
