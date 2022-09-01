package main

type RequestBody struct {
	User     string `json:"user"`
	Name     string `json:"name"`
	Document string `json:"document"`
	Role     string `json:"role"`
	Password string `json:"password"`
}
