package elasticsearch

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
)

func ConvertToElasticNote(note models.Note) models.ElasticNote {
	return models.ElasticNote{
		Id:            note.Id,
		Data:          string(note.Data),
		CreateTime:    note.CreateTime,
		UpdateTime:    note.UpdateTime,
		OwnerId:       note.OwnerId,
		Parent:        note.Parent,
		Children:      note.Children,
		Tags:          note.Tags,
		Collaborators: note.Collaborators,
	}
}

func ConvertToUsualNote(note models.ElasticNote) models.Note {
	return models.Note{
		Id:            note.Id,
		Data:          []byte(note.Data),
		CreateTime:    note.CreateTime,
		UpdateTime:    note.UpdateTime,
		OwnerId:       note.OwnerId,
		Parent:        note.Parent,
		Children:      note.Children,
		Tags:          note.Tags,
		Collaborators: note.Collaborators,
	}
}
