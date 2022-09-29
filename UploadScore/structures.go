package main

type RequestBody struct {
	Author   string `json:"author"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	Score    int    `json:"score"`
	Comments string `json:"comments"`
}
