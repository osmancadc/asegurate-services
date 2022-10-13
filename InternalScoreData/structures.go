package main

type InsertBody struct {
	Author    string `json:"author"`
	Objective string `json:"objective"`
	Score     int    `json:"score"`
	Comments  string `json:"comments"`
}

type UpdateBody struct {
	Document   string `json:"document"`
	Score      int    `json:"score"`
	Reputation int    `json:"reputation"`
}

type GetBody struct {
	Document string `json:"document"`
}

type RequestBody struct {
	Action     string     `json:"action"`
	InsertBody InsertBody `json:"insert_data,omitempty"`
	UpdateBody UpdateBody `json:"update_data,omitempty"`
	GetBody    GetBody    `json:"get_data,omitempty"`
}

type InternalScore struct {
	Score          float32
	PositiveScores int
	NegativeScores int
	Average60Days  float32
}
