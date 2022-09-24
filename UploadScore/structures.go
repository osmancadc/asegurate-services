package main

type RequestBody struct {
	Author   int    `json:"author"`
	Document string `json:"document"`
	Score    int    `json:"score"`
	Comments string `json:"comments"`
}
