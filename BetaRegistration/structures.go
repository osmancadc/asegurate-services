package main

type RequestBody struct {
	Name            string `json:"name"`
	Age             string `json:"age"`
	Cellphone       string `json:"cellphone"`
	Email           string `json:"email"`
	Smartphone      bool   `json:"smartphone"`
	OperativeSystem string `json:"os"`
	Commentary      string `json:"commentary"`
}
