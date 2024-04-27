package usecase

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/survey"
)

type SurveyUsecase struct {
	repo survey.SurveyRepo
}

func CreateSurveyUsecase(repo survey.SurveyRepo) *SurveyUsecase {
	return &SurveyUsecase{
		repo: repo,
	}
}
