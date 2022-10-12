/*
 Copyright 2018 Padduck, LLC
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
 	http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/comms"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Node struct {
	DB *gorm.DB
}

var nodeConnections = make(map[uint]*websocket.Conn)
var nodeLocker sync.Mutex

var LocalNode = &models.Node{
	ID:          0,
	Name:        "LocalNode",
	PublicHost:  "127.0.0.1",
	PrivateHost: "127.0.0.1",
	PublicPort:  8080,
	PrivatePort: 8080,
	SFTPPort:    5657,
	Secret:      strings.Replace(uuid.NewV4().String(), "-", "", -1),
}

func init() {
	nodeHost := config.WebHost.Value()
	sftpHost := config.SftpHost.Value()
	hostParts := strings.SplitN(nodeHost, ":", 2)
	sftpParts := strings.SplitN(sftpHost, ":", 2)

	if len(hostParts) == 2 {
		port, err := strconv.Atoi(hostParts[1])
		if err == nil {
			LocalNode.PublicPort = uint16(port)
			LocalNode.PrivatePort = uint16(port)
		}
	}
	if len(sftpParts) == 2 {
		port, err := strconv.Atoi(sftpParts[1])
		if err == nil {
			LocalNode.SFTPPort = uint16(port)
		}
	}
}

func (ns *Node) GetAll() ([]*models.Node, error) {
	var nodes []*models.Node

	res := ns.DB.Find(nodes)

	if res.Error != nil {
		return nil, res.Error
	}

	if config.PanelEnabled.Value() {
		hasLocal := false
		for _, v := range nodes {
			if v.IsLocal() {
				hasLocal = true
				break
			}
		}

		if !hasLocal {
			nodes = append(nodes, LocalNode)
		}
	}

	return nodes, nil
}

func (ns *Node) Get(id uint) (*models.Node, error) {
	model := &models.Node{}

	if id == LocalNode.ID && config.PanelEnabled.Value() {
		return LocalNode, nil
	}

	res := ns.DB.First(model, id)
	return model, res.Error
}

func (ns *Node) Update(model *models.Node) error {
	if model.ID == LocalNode.ID && config.PanelEnabled.Value() {
		return nil
	}

	res := ns.DB.Save(model)
	return res.Error
}

func (ns *Node) Delete(id uint) error {
	if id == LocalNode.ID && config.PanelEnabled.Value() {
		return errors.New("cannot delete local node")
	}

	model := &models.Node{
		ID: id,
	}

	var count int64
	ns.DB.Model(&models.Server{}).Where("node_id = ?", model.ID).Count(&count)
	if count > 0 {
		return pufferpanel.ErrNodeHasServers
	}

	res := ns.DB.Delete(model)
	return res.Error
}

func (ns *Node) Create(node *models.Node) error {
	res := ns.DB.Create(node)
	return res.Error
}

func (ns *Node) AddNodeConnection(node *models.Node, conn *websocket.Conn) {
	nodeLocker.Lock()
	defer nodeLocker.Unlock()

	//if we have a connection, terminate it
	existing := nodeConnections[node.ID]
	if existing != nil {
		_ = existing.Close()
	}

	go func(c *websocket.Conn) {
		messageType, d, err := c.ReadMessage()
		if err != nil {
			return
		}
		switch messageType {
		case websocket.PingMessage:
			{
				_ = c.WriteMessage(websocket.PongMessage, []byte{})
			}
		case websocket.BinaryMessage:
			{
				var msg comms.Message
				err = json.Unmarshal(d, &msg)
				if err != nil {
					return
				}

				ch, exists := pendingResponse.LoadAndDelete(msg.Id())
				if exists {
					//send message to the waiting channel
					ch2 := ch.(chan comms.Message)
					ch2 <- msg
				}

				//messageQueue.Store(msg.Id(), queueEntry{Time: time.Now(), Message: msg})
			}
		}
	}(conn)

	nodeConnections[node.ID] = conn
}

func (ns *Node) Send(node *models.Node, data interface{}) (comms.Message, error) {
	var id string

	if y, ok := data.(idInterface); ok {
		id = y.Id()
	} else if y, ok := data.(map[string]interface{}); ok {
		id = y["id"].(string)
	} else {
		return nil, errors.New("invalid message type")
	}

	ch := make(chan comms.Message, 1)

	d, _ := json.Marshal(data)
	pendingResponse.Store(id, ch)
	defer pendingResponse.Delete(id)

	err := ns.GetNodeConnection(node.ID).WriteMessage(websocket.BinaryMessage, d)
	if err != nil {
		return nil, err
	}

	//wait 30 seconds for a response back from the node
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var res comms.Message

	//either wait for timeout, or wait for response
	select {
	case <-ctx.Done():
		{
			err = ctx.Err()
		}
	case u := <-ch:
		{
			res = u
		}
	}

	return res, err
}

var pendingResponse = sync.Map{}

func (ns *Node) GetNodeConnection(id uint) *websocket.Conn {
	nodeLocker.Lock()
	defer nodeLocker.Unlock()

	return nodeConnections[id]
}

type idInterface interface {
	Id() string
}
