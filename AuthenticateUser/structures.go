package main

type RequestBody struct {
	Document string `json:"document"`
	Password string `json:"password"`
}

type User struct {
	UserId string
	Name   string
	Role   string
}
