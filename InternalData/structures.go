package main

type RequestBody struct {
	Action             string             `json:"action"`
	UpdateScoreBody    UpdateScoreBody    `json:"update_score_body,omitempty"`
	InsertScoreBody    InsertScoreBody    `json:"insert_score_body,omitempty"`
	InsertUserBody     InsertUserBody     `json:"insert_user_body,omitempty"`
	InsertPersonBody   InsertPersonBody   `json:"insert_person_body,omitempty"`
	GetScoreBody       GetScoreBody       `json:"get_score_body,omitempty"`
	GetUserByPhoneBody GetUserByPhoneBody `json:"get_user_by_phone_body,omitempty"`
	GetByDocumentBody  GetByDocumentBody  `json:"get_by_document_body,omitempty"`
}

// Structs Update

type UpdateScoreBody struct {
	Document   string `json:"document"`
	Score      int    `json:"score"`
	Reputation int    `json:"reputation"`
}

// Structs Insert

type InternalScore struct {
	Score          float32
	PositiveScores int
	NegativeScores int
	Average60Days  float32
}

type InsertScoreBody struct {
	Author    string `json:"author"`
	Objective string `json:"objective"`
	Score     int    `json:"score"`
	Comments  string `json:"comments"`
}

type InsertPersonBody struct {
	Document string `json:"document"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Gender   string `json:"gender"`
}

type InsertUserBody struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Document string `json:"document"`
	Role     string `json:"role"`
}

// Structs Get

type GetScoreBody struct {
	Document string `json:"document"`
}

type GetUserByPhoneBody struct {
	Phone string `json:"phone"`
}

type GetByDocumentBody struct {
	Document string `json:"document"`
}
