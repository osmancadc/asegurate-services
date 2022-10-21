package main

type RequestBody struct {
	Action            string            `json:"action"`
	PersonBody        Person            `json:"person_body,omitempty"`
	UserBody          User              `json:"user_body,omitempty"`
	ScoreBody         ScoreBody         `json:"score_body,omitempty"`
	GetByPhoneBody    GetByPhoneBody    `json:"get_by_phone_body,omitempty"`
	GetByDocumentBody GetByDocumentBody `json:"get_by_document_body,omitempty"`
}

type InternalScore struct {
	Score          float32
	PositiveScores int
	NegativeScores int
	Average60Days  float32
}

type ScoreBody struct {
	Author    string `json:"author"`
	Objective string `json:"objective"`
	Score     int    `json:"score"`
	Comments  string `json:"comments"`
}

type GetByPhoneBody struct {
	Phone  string   `json:"phone"`
	Fields []string `json:"fields,omitempty"`
}

type GetByDocumentBody struct {
	Document string   `json:"document"`
	Fields   []string `json:"fields,omitempty"`
}

type Account struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Photo    string `json:"photo,omitempty"`
	Gender   string `json:"gender,omitempty"`
	Document string `json:"document"`
	Role     string `json:"role,omitempty"`
}

// Database Models

type User struct {
	UserId   int    `gorm:"<-:false" json:"user_id,omitempty"`
	Document string `gorm:"<-:create" json:"document,omitempty"`
	Email    string `gorm:"<-" json:"email,omitempty"`
	Phone    string `gorm:"<-" json:"phone,omitempty"`
	Password string `gorm:"<-" json:"password,omitempty"`
	Role     string `gorm:"<-" json:"role,omitempty"`
}

type Person struct {
	Document   string `gorm:"<-:create" json:"document,omitempty"`
	Name       string `gorm:"<-:create" json:"name,omitempty"`
	Lastname   string `gorm:"<-:create" json:"lastname,omitempty"`
	Gender     string `gorm:"<-" json:"gender,omitempty"`
	Reputation int    `gorm:"<-update" json:"reputation,omitempty"`
	Photo      string `gorm:"<-" json:"photo,omitempty"`
	LastUpdate string `gorm:"<-:update" json:"last_update,omitempty"`
}

type Score struct {
	ID           int    `gorm:"<-:create"`
	Author       int    `gorm:"<-:create"`
	Objective    string `gorm:"<-:create"`
	Score        int    `gorm:"<-:create"`
	Comments     string `gorm:"<-:create"`
	CreationDate string `gorm:"<-:false"`
}
