package main

type RequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	UserId string
	Name   string
	Role   string
}
