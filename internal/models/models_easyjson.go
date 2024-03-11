// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels(in *jlexer.Lexer, out *UserFormData) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "username":
			out.Username = string(in.String())
		case "password":
			out.Password = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels(out *jwriter.Writer, in UserFormData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix[1:])
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserFormData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserFormData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserFormData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserFormData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels1(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.Id).UnmarshalText(data))
			}
		case "description":
			out.Description = string(in.String())
		case "username":
			out.Username = string(in.String())
		case "create_time":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreateTime).UnmarshalJSON(data))
			}
		case "image_path":
			out.ImagePath = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels1(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.RawText((in.Id).MarshalText())
	}
	if in.Description != "" {
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"create_time\":"
		out.RawString(prefix)
		out.Raw((in.CreateTime).MarshalJSON())
	}
	{
		const prefix string = ",\"image_path\":"
		out.RawString(prefix)
		out.String(string(in.ImagePath))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels1(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels2(in *jlexer.Lexer, out *ProfileUpdatePayload) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "description":
			out.Description = string(in.String())
		case "password":
			easyjsonD2b7633eDecode(in, &out.Password)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels2(out *jwriter.Writer, in ProfileUpdatePayload) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Description != "" {
		const prefix string = ",\"description\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Description))
	}
	if true {
		const prefix string = ",\"password\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		easyjsonD2b7633eEncode(out, in.Password)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ProfileUpdatePayload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ProfileUpdatePayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ProfileUpdatePayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ProfileUpdatePayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels2(l, v)
}
func easyjsonD2b7633eDecode(in *jlexer.Lexer, out *struct {
	Old string `json:"old"`
	New string `json:"new"`
}) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "old":
			out.Old = string(in.String())
		case "new":
			out.New = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncode(out *jwriter.Writer, in struct {
	Old string `json:"old"`
	New string `json:"new"`
}) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"old\":"
		out.RawString(prefix[1:])
		out.String(string(in.Old))
	}
	{
		const prefix string = ",\"new\":"
		out.RawString(prefix)
		out.String(string(in.New))
	}
	out.RawByte('}')
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels3(in *jlexer.Lexer, out *Note) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.Id).UnmarshalText(data))
			}
		case "data":
			if in.IsNull() {
				in.Skip()
				out.Data = nil
			} else {
				out.Data = in.Bytes()
			}
		case "create_time":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreateTime).UnmarshalJSON(data))
			}
		case "update_time":
			if in.IsNull() {
				in.Skip()
				out.UpdateTime = nil
			} else {
				if out.UpdateTime == nil {
					out.UpdateTime = new(time.Time)
				}
				if data := in.Raw(); in.Ok() {
					in.AddError((*out.UpdateTime).UnmarshalJSON(data))
				}
			}
		case "owner_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.OwnerId).UnmarshalText(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels3(out *jwriter.Writer, in Note) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.RawText((in.Id).MarshalText())
	}
	if len(in.Data) != 0 {
		const prefix string = ",\"data\":"
		out.RawString(prefix)
		out.Base64Bytes(in.Data)
	}
	{
		const prefix string = ",\"create_time\":"
		out.RawString(prefix)
		out.Raw((in.CreateTime).MarshalJSON())
	}
	if in.UpdateTime != nil {
		const prefix string = ",\"update_time\":"
		out.RawString(prefix)
		out.Raw((*in.UpdateTime).MarshalJSON())
	}
	{
		const prefix string = ",\"owner_id\":"
		out.RawString(prefix)
		out.RawText((in.OwnerId).MarshalText())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Note) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Note) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Note) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Note) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels3(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels4(in *jlexer.Lexer, out *JwtPayload) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.Id).UnmarshalText(data))
			}
		case "Username":
			out.Username = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels4(out *jwriter.Writer, in JwtPayload) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Id\":"
		out.RawString(prefix[1:])
		out.RawText((in.Id).MarshalText())
	}
	{
		const prefix string = ",\"Username\":"
		out.RawString(prefix)
		out.String(string(in.Username))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v JwtPayload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v JwtPayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *JwtPayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *JwtPayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels4(l, v)
}
