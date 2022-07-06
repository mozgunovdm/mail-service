package models

import "net/http"

type Response struct{}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Task struct {
	Id          string   `json:"id" bson:"_id"`
	Title       string   `json:"title" bson:"title"`
	Description string   `json:"description" bson:"description,omitempty"`
	MailList    []string `json:"mailList" bson:"mailList"`
	AuthorId    string   `json:"authorId" bson:"authId"`
	Status      string   `json:"status" bson:"status"`
	Response
}

type Body struct {
	Id          string   `json:"id" bson:"_id"`
	Title       string   `json:"title" bson:"title"`
	Description string   `json:"description" bson:"description,omitempty"`
	MailList    []string `json:"mailList" bson:"mailList"`
	AuthorId    string   `json:"authorId" bson:"status"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
	Response
}
