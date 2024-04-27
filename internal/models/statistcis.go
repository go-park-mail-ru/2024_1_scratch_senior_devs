package models

import "github.com/satori/uuid"

type Question struct {
	Id           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	QuestionType string    `json:"question_type"`
	SurveyId     uuid.UUID `json:"survey_id"`
}

type Result struct {
	Id         uuid.UUID `json:"id"`
	QuestionId uuid.UUID `json:"question_id"`
	Voice      int       `json:"voice"`
}

type Vote struct {
	QuestionId uuid.UUID `json:"question_id"`
	Voice      int       `json:"voice"`
}
