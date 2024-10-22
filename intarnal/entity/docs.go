package entity

import "time"

type Meta struct {
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	File    bool      `json:"file"`
	Public  bool      `json:"public"`
	Mime    string    `json:"mime"`
	Grant   []string  `json:"grant"`
	Created time.Time `json:"created"`
}

type MetaSlice []Meta

func (m MetaSlice) Len() int {
	return len(m)
}

func (m MetaSlice) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m MetaSlice) Less(i, j int) bool {
	if m[i].Name == m[j].Name {
		return m[i].Created.Before(m[j].Created)
	}
	return m[i].Name < m[j].Name
}

func NewMeta(name string, file, public bool, mime string, grant string) *Meta {
	return &Meta{
		Name:   name,
		File:   file,
		Public: public,
		Mime:   mime,
		Grant:  []string{grant},
	}
}

type Document struct {
	Body     []byte `json:"body"`
	Mime     string `json:"mime"`
	Name     string `json:"name"`
	JsonBody string `json:"json_body"`
}

func NewDocument(body []byte, mime, name, jsonBody string) *Document {
	return &Document{
		Body:     body,
		Mime:     mime,
		Name:     name,
		JsonBody: jsonBody,
	}
}

type SearchDocuments struct {
	Login  string `json:"login"`
	Limit  int    `json:"limit"`
	Name   string `json:"name"`
	File   *bool  `json:"file"`
	Public *bool  `json:"public"`
	Mime   string `json:"mime"`
}
