package comms

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

func Send(conn *websocket.Conn, data interface{}) error {
	d, _ := json.Marshal(data)
	return conn.WriteMessage(websocket.BinaryMessage, d)
}
