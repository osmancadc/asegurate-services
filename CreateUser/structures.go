package main

type RequestBody struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Document string `json:"document"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Password string `json:"password"`
}
