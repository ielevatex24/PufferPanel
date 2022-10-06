package comms

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferpanel/v3/comms"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/daemon/programs"
	"net/http"
	"strings"
	"sync"
	"time"
)

var _conn *websocket.Conn
var locker sync.Mutex

func StartConnection() {
	go listen()
}

func recoverConnection() {
	locker.Lock()
	defer locker.Unlock()

	if _conn != nil {
		//if no error sending a ping, then we have a connection
		if _conn.WriteMessage(websocket.PingMessage, []byte{}) == nil {
			return
		}
	}

	baseUrl := config.PanelUrl.Value()
	secret := config.DaemonSecret.Value()

	if baseUrl == "" && config.PanelEnabled.Value() {
		baseUrl = strings.Replace(config.PanelUrl.Value(), "0.0.0.0:", "127.0.0.1:", 1)
	}

	header := http.Header{}
	header.Set("Authorization", "Node "+secret)
	conn, _, err := websocket.DefaultDialer.Dial(baseUrl, header)

	for err != nil {
		time.Sleep(time.Second * 5)
		conn, _, err = websocket.DefaultDialer.Dial(baseUrl, header)
	}

	if conn != nil {
		_conn = conn
	}
}

func listen() {
	for {
		messageType, d, err := _conn.ReadMessage()
		if err != nil {
			recoverConnection()
		}
		switch messageType {
		case websocket.PingMessage:
			{
				_ = _conn.WriteMessage(websocket.PongMessage, []byte{})
			}
		case websocket.BinaryMessage:
			{
				//process this request in a new routine, so we don't stall other processing
				go func(processData []byte) {
					var msg comms.Message
					err = json.Unmarshal(d, &msg)
					if err != nil {
						return
					}

					//send confirmation we got the request, and are working on it
					_ = Send(comms.NewConfirmation(msg.Id()))
					switch msg.Type() {
					case comms.StartServerType():
						{
							data := comms.Cast[comms.StartServer](msg)

							prg, err := programs.Get(data.Server)
							if err != nil {
								return
							}
							err = prg.Start()
							if err != nil {
								return
							}
						}
					}
				}(d)
			}
		}
	}
}

func Send(data interface{}) error {
	if _conn == nil {
		recoverConnection()
	}

	d, _ := json.Marshal(data)

	//we will attempt to send the message twice, once if we think we have a connection, and then again once we've
	//actually believe we have one
	err := _conn.WriteMessage(websocket.BinaryMessage, d)
	if err != nil {
		recoverConnection()
		return _conn.WriteMessage(websocket.BinaryMessage, d)
	}
	return nil
}
