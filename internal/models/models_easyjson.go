// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels(in *jlexer.Lexer, out *passwords) {
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels(out *jwriter.Writer, in passwords) {
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

// MarshalJSON supports json.Marshaler interface
func (v passwords) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v passwords) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *passwords) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *passwords) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels1(in *jlexer.Lexer, out *UserFormData) {
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
		case "code":
			out.Code = string(in.String())
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels1(out *jwriter.Writer, in UserFormData) {
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
	if in.Code != "" {
		const prefix string = ",\"code\":"
		out.RawString(prefix)
		out.String(string(in.Code))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserFormData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserFormData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserFormData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserFormData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels1(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels2(in *jlexer.Lexer, out *UserForSwagger) {
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
		case "second_factor":
			out.SecondFactor = bool(in.Bool())
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels2(out *jwriter.Writer, in UserForSwagger) {
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
	{
		const prefix string = ",\"second_factor\":"
		out.RawString(prefix)
		out.Bool(bool(in.SecondFactor))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserForSwagger) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserForSwagger) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserForSwagger) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserForSwagger) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels2(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels3(in *jlexer.Lexer, out *User) {
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
		case "second_factor":
			out.SecondFactor = Secret(in.String())
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels3(out *jwriter.Writer, in User) {
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
	{
		const prefix string = ",\"second_factor\":"
		out.RawString(prefix)
		out.Raw((in.SecondFactor).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels3(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels4(in *jlexer.Lexer, out *UpsertNoteRequestForSwagger) {
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
		case "data":
			(out.Data).UnmarshalEasyJSON(in)
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels4(out *jwriter.Writer, in UpsertNoteRequestForSwagger) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix[1:])
		(in.Data).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UpsertNoteRequestForSwagger) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UpsertNoteRequestForSwagger) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UpsertNoteRequestForSwagger) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UpsertNoteRequestForSwagger) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels4(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels5(in *jlexer.Lexer, out *UpsertNoteRequest) {
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
		case "data":
			if m, ok := out.Data.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Data.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Data = in.Interface()
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels5(out *jwriter.Writer, in UpsertNoteRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix[1:])
		if m, ok := in.Data.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Data.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Data))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UpsertNoteRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UpsertNoteRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UpsertNoteRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UpsertNoteRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels5(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels6(in *jlexer.Lexer, out *SignUpPayloadForSwagger) {
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels6(out *jwriter.Writer, in SignUpPayloadForSwagger) {
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
func (v SignUpPayloadForSwagger) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SignUpPayloadForSwagger) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SignUpPayloadForSwagger) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SignUpPayloadForSwagger) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels6(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels7(in *jlexer.Lexer, out *ProfileUpdatePayload) {
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
			(out.Password).UnmarshalEasyJSON(in)
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels7(out *jwriter.Writer, in ProfileUpdatePayload) {
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
		(in.Password).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ProfileUpdatePayload) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ProfileUpdatePayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ProfileUpdatePayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ProfileUpdatePayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels7(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels8(in *jlexer.Lexer, out *NoteForSwagger) {
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
			(out.Data).UnmarshalEasyJSON(in)
		case "create_time":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreateTime).UnmarshalJSON(data))
			}
		case "update_time":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.UpdateTime).UnmarshalJSON(data))
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels8(out *jwriter.Writer, in NoteForSwagger) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.RawText((in.Id).MarshalText())
	}
	if true {
		const prefix string = ",\"data\":"
		out.RawString(prefix)
		(in.Data).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"create_time\":"
		out.RawString(prefix)
		out.Raw((in.CreateTime).MarshalJSON())
	}
	if true {
		const prefix string = ",\"update_time\":"
		out.RawString(prefix)
		out.Raw((in.UpdateTime).MarshalJSON())
	}
	{
		const prefix string = ",\"owner_id\":"
		out.RawString(prefix)
		out.RawText((in.OwnerId).MarshalText())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v NoteForSwagger) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NoteForSwagger) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NoteForSwagger) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NoteForSwagger) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels8(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels9(in *jlexer.Lexer, out *NoteDataForSwagger) {
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
		case "title":
			out.Title = string(in.String())
		case "content":
			out.Content = string(in.String())
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels9(out *jwriter.Writer, in NoteDataForSwagger) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix[1:])
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"content\":"
		out.RawString(prefix)
		out.String(string(in.Content))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v NoteDataForSwagger) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NoteDataForSwagger) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NoteDataForSwagger) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NoteDataForSwagger) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels9(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels10(in *jlexer.Lexer, out *Note) {
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
			if data := in.Raw(); in.Ok() {
				in.AddError((out.UpdateTime).UnmarshalJSON(data))
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels10(out *jwriter.Writer, in Note) {
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
	{
		const prefix string = ",\"update_time\":"
		out.RawString(prefix)
		out.Raw((in.UpdateTime).MarshalJSON())
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
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels10(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Note) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels10(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Note) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels10(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Note) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels10(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels11(in *jlexer.Lexer, out *JwtPayload) {
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels11(out *jwriter.Writer, in JwtPayload) {
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
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels11(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v JwtPayload) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels11(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *JwtPayload) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels11(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *JwtPayload) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels11(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels12(in *jlexer.Lexer, out *ElasticNoteUpdate) {
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
		case "doc":
			(out.Doc).UnmarshalEasyJSON(in)
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels12(out *jwriter.Writer, in ElasticNoteUpdate) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"doc\":"
		out.RawString(prefix[1:])
		(in.Doc).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ElasticNoteUpdate) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels12(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ElasticNoteUpdate) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels12(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ElasticNoteUpdate) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels12(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ElasticNoteUpdate) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels12(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels13(in *jlexer.Lexer, out *ElasticNote) {
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
			out.Data = string(in.String())
		case "elastic_data":
			if m, ok := out.ElasticData.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.ElasticData.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.ElasticData = in.Interface()
			}
		case "create_time":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreateTime).UnmarshalJSON(data))
			}
		case "update_time":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.UpdateTime).UnmarshalJSON(data))
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels13(out *jwriter.Writer, in ElasticNote) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.RawText((in.Id).MarshalText())
	}
	if in.Data != "" {
		const prefix string = ",\"data\":"
		out.RawString(prefix)
		out.String(string(in.Data))
	}
	if in.ElasticData != nil {
		const prefix string = ",\"elastic_data\":"
		out.RawString(prefix)
		if m, ok := in.ElasticData.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.ElasticData.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.ElasticData))
		}
	}
	{
		const prefix string = ",\"create_time\":"
		out.RawString(prefix)
		out.Raw((in.CreateTime).MarshalJSON())
	}
	{
		const prefix string = ",\"update_time\":"
		out.RawString(prefix)
		out.Raw((in.UpdateTime).MarshalJSON())
	}
	{
		const prefix string = ",\"owner_id\":"
		out.RawString(prefix)
		out.RawText((in.OwnerId).MarshalText())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ElasticNote) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels13(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ElasticNote) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels13(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ElasticNote) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels13(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ElasticNote) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels13(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels14(in *jlexer.Lexer, out *Attach) {
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
		case "path":
			out.Path = string(in.String())
		case "note_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.NoteId).UnmarshalText(data))
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels14(out *jwriter.Writer, in Attach) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.RawText((in.Id).MarshalText())
	}
	{
		const prefix string = ",\"path\":"
		out.RawString(prefix)
		out.String(string(in.Path))
	}
	{
		const prefix string = ",\"note_id\":"
		out.RawString(prefix)
		out.RawText((in.NoteId).MarshalText())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Attach) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels14(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Attach) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels14(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Attach) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels14(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Attach) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20241ScratchSeniorDevsInternalModels14(l, v)
}
