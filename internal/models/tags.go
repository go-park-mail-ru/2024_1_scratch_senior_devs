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

type UpdateTagRequest struct {
	OldTag string `json:"old_name"`
	NewTag string `json:"new_name"`
}

func (payload *UpdateTagRequest) Sanitize() {
	payload.NewTag = html.EscapeString(payload.NewTag)
	payload.OldTag = html.EscapeString(payload.OldTag)

}

func (payload *UpdateTagRequest) Validate() error {
	if utf8.RuneCountInString(payload.NewTag) > maxTagLength {
		return errors.New("tag name too long")
	}
	if utf8.RuneCountInString(payload.OldTag) > maxTagLength {
		return errors.New("tag name too long")
	}
	return nil
}

type GetTagsResponse struct {
	Tags []string `json:"tags"`
}
