package main

type RequestBody struct {
	User     string `json:"user"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Password string `json:"password"`
}
