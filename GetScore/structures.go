package main

import "time"

type RequestBody struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type StoredReputation struct {
	Name       string `json:"name,omitempty"`
	Lastname   string `json:"lastname,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Reputation int    `json:"reputation,omitempty"`
	LastUpdate string `json:"last_update,omitempty"`
	Photo      string `json:"photo,omitempty"`
}

type InternalScore struct {
	Score          float32 `json:"score"`
	PositiveScores int     `json:"positive_scores"`
	NegativeScores int     `json:"negative_scores"`
	Average60Days  float32 `json:"average_60_days"`
}

type ExternalProccedings struct {
	FormalComplaints  int `json:"formal_complaints"`
	FormalRecentYear  int `json:"recent_complain_year"`
	Formal5YearsTotal int `json:"five_years_amount"`
}

type InvokePayload struct {
	Body string `json:"body"`
}

type InvokeBody struct {
	Action          string            `json:"action"`
	GetByDocument   GetByDocumentBody `json:"get_by_document_body,omitempty"`
	GetExternalBody GetExternalBody   `json:"name_body,omitempty"`
	GetByPhone      GetByPhoneBody    `json:"get_by_phone_body,omitempty"`
	Person          Person            `json:"person_body,omitempty"`
}

type GetByDocumentBody struct {
	Document string   `json:"document"`
	Fields   []string `josn:"fields"`
}

type GetByPhoneBody struct {
	Phone  string   `json:"phone"`
	Fields []string `json:"fields"`
}

type GetExternalBody struct {
	Document string `json:"document"`
	Type     string `json:"document_type"`
}

type PredictScoreBody struct {
	InternalScore          float32 `json:"internal_score"`
	InternalPositiveScores int     `json:"positive_score"`
	InternalNegativeScores int     `json:"negative_score"`
	InternalAverage60Days  float32 `json:"average_score"`
	MeliScore              int     `json:"meli_score"`
	MeliPositiveScores     int     `json:"meli_positive"`
	MeliNegativeScores     int     `json:"meli_negative"`
	MeliAmountSales        int     `json:"sales_completed"`
	FormalComplaints       int     `json:"formal_complaints"`
	FormalRecentYear       int     `json:"recent_complain_year"`
	Formal5YearsTotal      int     `json:"five_years_amount"`
}

type InvokeResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

type PredictionResponse struct {
	Score      int `json:"score"`
	Reputation int `json:"reputation"`
}

type Score struct {
	Document   string `json:"document"`
	Name       string `json:"name"`
	Gender     string `json:"gender"`
	Score      int    `json:"score"`
	Reputation int    `json:"reputation"`
	Stars      int    `json:"stars"`
	Updated    string `json:"last_updated,omitempty"`
	Certified  bool   `json:"certified"`
	Photo      string `json:"photo"`
}

type Person struct {
	Document   string    `json:"document,omitempty"`
	Name       string    `json:"name,omitempty"`
	Lastname   string    `json:"lastname,omitempty"`
	Gender     string    `json:"gender,omitempty"`
	Reputation int       `json:"reputation,omitempty"`
	Photo      string    `json:"photo,omitempty"`
	LastUpdate time.Time `json:"last_update,omitempty"`
}
