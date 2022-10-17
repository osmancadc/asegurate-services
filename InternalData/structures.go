package main

import "time"

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
	Phone string `json:"phone"`
}

type GetByDocumentBody struct {
	Document string `json:"document"`
}

// Database Models

type User struct {
	UserId   int    `gorm:"<-:false"`
	Document string `gorm:"<-:create"`
	Email    string `gorm:"<-"`
	Phone    string `gorm:"<-"`
	Password string `gorm:"<-"`
	Role     string `gorm:"<-"`
}

type Person struct {
	Document   string    `gorm:"<-:create" json:"document"`
	Name       string    `gorm:"<-:create" json:"name,omitempty"`
	Lastname   string    `gorm:"<-:create" json:"lastname,omitempty"`
	Gender     string    `gorm:"<-" json:"gender,omitempty"`
	Score      int       `gorm:"<-update" json:"score,omitempty"`
	Reputation int       `gorm:"<-update" json:"reputation,omitempty"`
	Photo      string    `gorm:"<-" json:"photo,omitempty"`
	LastUpdate time.Time `gorm:"<-:update"`
}

type Score struct {
	ID           int       `gorm:"<-:create"`
	Author       int       `gorm:"<-:create"`
	Objective    string    `gorm:"<-:create"`
	Score        int       `gorm:"<-:create"`
	Comments     string    `gorm:"<-:create"`
	CreationDate time.Time `gorm:"<-:false"`
}
