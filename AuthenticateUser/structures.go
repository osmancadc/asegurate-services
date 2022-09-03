package main

type RequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	UserId int
	Name   string
	Role   string
}
