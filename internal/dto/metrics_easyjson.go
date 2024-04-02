/*
 * This file was last modified at 2024-03-18 23:47 by Victor N. Skurikhin.
 * metrics_easyjson.go
 * $Id$
 */

// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package dto

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

func easyjson2220f231DecodeGithubComVskurikhinGometricsInternalDto(in *jlexer.Lexer, out *Metrics) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Metrics, 0, 1)
			} else {
				*out = Metrics{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Metric
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2220f231EncodeGithubComVskurikhinGometricsInternalDto(out *jwriter.Writer, in Metrics) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Metrics) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2220f231EncodeGithubComVskurikhinGometricsInternalDto(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Metrics) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2220f231EncodeGithubComVskurikhinGometricsInternalDto(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Metrics) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2220f231DecodeGithubComVskurikhinGometricsInternalDto(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Metrics) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2220f231DecodeGithubComVskurikhinGometricsInternalDto(l, v)
}
