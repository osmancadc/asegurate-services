package main

type User struct {
	Document string `json:"document"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Photo    string `json:"photo"`
}
