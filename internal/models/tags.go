package models

import (
	"errors"
	"html"
	"unicode/utf8"
)

const (
	maxTagLength = 20
)

type TagRequest struct {
	TagName string `json:"tag_name"`
}

func (payload *TagRequest) Sanitize() {
	payload.TagName = html.EscapeString(payload.TagName)

}

func (payload *TagRequest) Validate() error {
	if utf8.RuneCountInString(payload.TagName) > maxTagLength {
		return errors.New("tag name too long")
	}
	return nil
}
