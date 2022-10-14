package main

type RequestBody struct {
	Action          string          `json:"action"`
	InsertScoreBody InsertScoreBody `json:"insert_score_body,omitempty"`
	UpdateScoreBody UpdateScoreBody `json:"update_score_body,omitempty"`
	GetScoreBody    GetScoreBody    `json:"get_score_body,omitempty"`
	GetByPhoneBody  GetByPhoneBody  `json:"get_by_phone_body,omitempty"`
}

type InsertScoreBody struct {
	Author    string `json:"author"`
	Objective string `json:"objective"`
	Score     int    `json:"score"`
	Comments  string `json:"comments"`
}

type UpdateScoreBody struct {
	Document   string `json:"document"`
	Score      int    `json:"score"`
	Reputation int    `json:"reputation"`
}

type GetScoreBody struct {
	Document string `json:"document"`
}

type GetByPhoneBody struct {
	Phone string `json:"phone"`
}

type InternalScore struct {
	Score          float32
	PositiveScores int
	NegativeScores int
	Average60Days  float32
}
