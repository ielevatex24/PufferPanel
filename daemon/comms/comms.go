package comms

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferpanel/v3/comms"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/daemon/programs"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/models"
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

	if config.PanelEnabled.Value() {
		secret = models.LocalNode.Secret
	}

	if baseUrl == "" && config.PanelEnabled.Value() {
		baseUrl = strings.Replace(config.WebHost.Value(), "0.0.0.0:", "127.0.0.1:", 1) + "/auth/node/socket"
	}
	baseUrl = strings.Replace(baseUrl, "https://", "wss://", 1)
	baseUrl = strings.Replace(baseUrl, "http://", "ws://", 1)
	if !strings.HasPrefix(baseUrl, "ws://") && !strings.HasPrefix(baseUrl, "wss://") {
		baseUrl = "ws://" + baseUrl
	}
	if !strings.HasSuffix(baseUrl, "/auth/node/socket") {
		baseUrl = baseUrl + "/auth/node/socket"
	}

	header := http.Header{}
	header.Set("Authorization", "Node "+secret)

	logging.Debug.Printf("opening websocket connection to %s", baseUrl)
	conn, _, err := websocket.DefaultDialer.Dial(baseUrl, header)

	for err != nil {
		logging.Error.Printf("error re-establishing socket connection to panel: %s", err)
		time.Sleep(time.Second * 5)
		conn, _, err = websocket.DefaultDialer.Dial(baseUrl, header)
	}

	logging.Info.Printf("re-established socket connection to panel")

	if conn != nil {
		_conn = conn
	}
}

func listen() {
	defer func() {
		if err := recover(); err != nil {
			logging.Error.Printf("critical error re-establishing socket connection to panel: %s", err)
		}
	}()
	if _conn == nil {
		recoverConnection()
	}
	for {
		messageType, d, e := _conn.ReadMessage()
		if e != nil {
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
					err := json.Unmarshal(d, &msg)
					if err != nil {
						return
					}

					switch msg.Type() {
					case comms.ConfirmationType():
						break
					case comms.StartServerType():
						{
							data := comms.Cast[comms.StartServer](msg)

							prg, err := programs.Get(data.Server)
							if err != nil {
								_ = Send(comms.NewError(msg.Id(), err))
								return
							}
							err = prg.Start()
							if err != nil {
								_ = Send(comms.NewError(msg.Id(), err))
								return
							}
							_ = Send(comms.NewConfirmation(msg.Id()))
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
	//actually believed we have one
	err := _conn.WriteMessage(websocket.BinaryMessage, d)
	if err != nil {
		recoverConnection()
		return _conn.WriteMessage(websocket.BinaryMessage, d)
	}
	return nil
}
