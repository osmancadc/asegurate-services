package main

type RequestBody struct {
	Author    string `json:"author"`
	Type      string `json:"type,omitempty"`
	Objective string `json:"objective"`
	Score     int    `json:"score"`
	Comments  string `json:"comments"`
}

type FindByPhoneBody struct {
	Phone string `json:"phone"`
}

type InvokeBody struct {
	Action          string          `json:"action"`
	InsertData      RequestBody     `json:"insert_score_body,omitempty"`
	FindByPhoneData FindByPhoneBody `json:"get_by_phone_body,omitempty"`
}

type InvokePayload struct {
	Body string `json:"body"`
}

type InvokeResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type FindByPhoneResponseBody struct {
	Document string `json:"document"`
}
