package comms

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cast"
)

type Message map[string]interface{}

func (m Message) Type() string {
	return m.String("type")
}

func (m Message) Id() string {
	return m.String("id")
}

func (m Message) String(key string) string {
	return cast.ToString(m[key])
}

func Cast[T any](input Message) T {
	var data T
	temp, _ := json.Marshal(input)
	_ = json.Unmarshal(temp, &data)
	return data
}

func NewId() string {
	return uuid.NewV4().String()
}

type Confirmation struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

func ConfirmationType() string { return "confirmation" }
func NewConfirmation(id string) Confirmation {
	return Confirmation{
		Type: ConfirmationType(),
		Id:   id,
	}
}

type Error struct {
	Type   string `json:"type"`
	Id     string `json:"id"`
	Error  error  `json:"error"`
	Server string `json:"server"`
}

func ErrorType() string { return "error" }
func NewError(id string, err error) Error {
	return Error{
		Type:  ErrorType(),
		Id:    id,
		Error: err,
	}
}
func NewErrorOnServer(id string, err error, server string) Error {
	e := NewError(id, err)
	e.Server = server
	return e
}

func IsSuccess(msg Message) bool {
	return msg.Type() == ConfirmationType()
}

func IsError(msg Message) bool {
	return msg.Type() == ErrorType()
}
