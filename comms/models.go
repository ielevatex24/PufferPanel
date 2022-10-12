package comms

import (
	"encoding/json"
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

type Confirmation struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

func ConfirmationType() string { return "confirmation" }
func NewConfirmation(id string) Confirmation {
	return Confirmation{
		Type: StartServerType(),
		Id:   id,
	}
}

type StartServer struct {
	Type   string `json:"type"`
	Id     string `json:"id"`
	Server string `json:"server"`
}

func StartServerType() string { return "server start" }
func NewStartServer(server string) StartServer {
	return StartServer{
		Type:   StartServerType(),
		Server: server,
	}
}
