package validation

import (
	"encoding/json"
	"errors"
)

func CheckNoteTitle(noteData []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(noteData, &m)
	if err != nil {
		return err
	}

	title, ok := m["title"]
	if !ok {
		return errors.New("title not found")
	}

	stringTitle, ok := title.(string)
	if !ok {
		return errors.New("invalid title format")
	}

	if len(stringTitle) == 0 {
		return errors.New("empty title")
	}

	return nil
}
