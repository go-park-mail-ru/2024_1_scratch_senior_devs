package repo

import "github.com/jackc/pgtype/pgxtype"

type SurveyRepo struct {
	db pgxtype.Querier
}

func CreateSurveyRepo(db pgxtype.Querier) *SurveyRepo {
	return &SurveyRepo{
		db: db,
	}
}
