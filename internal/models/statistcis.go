package models

import "github.com/satori/uuid"

type Question struct {
	Id           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	MinMark      int       `json:"min_mark"`
	Skip         int       `json:"skip"`
	QuestionType string    `json:"question_type"`
	Number       int       `json:"number"`
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
type Stat struct {
	QuestionId   uuid.UUID `json:"question_id"`
	Title        string    `json:"title"`
	QuestionType string    `json:"question_type"`
	Voice        int       `json:"voice"`
	Count        int       `json:"count"`
	Type         string    `json:"-"`
}

type CreateQuestionRequest struct {
	Title        string `json:"title"`
	MinMark      int    `json:"min_mark"`
	Skip         int    `json:"skip"`
	QuestionType string `json:"question_type"`
}

type StatResponse struct {
	QuestionId   uuid.UUID   `json:"question_id"`
	Title        string      `json:"title"`
	QuestionType string      `json:"question_type"`
	Stats        interface{} `json:"stats"`
	Value        float64     `json:"value"`
}

type CreateSurveyRequest struct {
	Questions []CreateQuestionRequest `json:"questions"`
}
