package models

import "github.com/satori/uuid"

type Question struct {
	Id           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	QuestionType string    `json:"question_type"`
	Number       int       `json:"number"`
	SurveyId     uuid.UUID `json:"survey_id"`
}

type Result struct {
	Id         uuid.UUID `json:"id"`
	QuestionId uuid.UUID `json:"question_id"`
	Voice      int       `json:"voice"`
}
