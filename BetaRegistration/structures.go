package main

type RequestBody struct {
	Name            string `json:"name"`
	Age             int    `json:"age"`
	Cellphone       string `json:"cellphone"`
	Email           string `json:"email"`
	Associate       string `json:"associate"`
	Smartphone      bool   `json:"smartphone"`
	OperativeSystem string `json:"os"`
	Commentary      string `json:"commentary"`
}
