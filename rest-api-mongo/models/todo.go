package models

import (
	"encoding/json"
	"io"
)

// TodoModel
type Todo struct {
	ID     int64  `json:"id" bson:"_id, omitempty"`
	Name   string `json:"name" bson:"name, omitempty"`
	Author string `json:"author" bson:"author, omitempty"`
}

// FromJSON for unmarshalling/decoding TodoObject from payload
func (t *Todo) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}

// ToJSON for marshalling TodoObject and sending it in the response
func (t *Todo) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

// TodoRepository interface => All available methods for TodoCollection
type TodoRepository interface {
	GetAll() (*[]Todo, error)
	Get(ID int64) (*Todo, error)
	Save(todo *Todo) error
	Delete(ID int64) error
}
