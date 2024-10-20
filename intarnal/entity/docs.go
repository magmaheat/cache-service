package entity

import "time"

type Meta struct {
	Id      string    `json:"token"`
	Name    string    `json:"name"`
	File    bool      `json:"file"`
	Public  bool      `json:"public"`
	Mime    string    `json:"mime"`
	Grant   []string  `json:"grant"`
	Created time.Time `json:"created"`
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
