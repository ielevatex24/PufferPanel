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
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Node struct {
	DB *gorm.DB
}

var nodeConnections = make(map[uint]*websocket.Conn)
var nodeLocker sync.Mutex

func (ns *Node) GetAll() (*models.Nodes, error) {
	nodes := &models.Nodes{}

	res := ns.DB.Find(nodes)

	return nodes, res.Error
}

func (ns *Node) Get(id uint) (*models.Node, error) {
	model := &models.Node{}

	res := ns.DB.First(model, id)

	return model, res.Error
}

func (ns *Node) Update(model *models.Node) error {
	res := ns.DB.Save(model)
	return res.Error
}

func (ns *Node) Delete(id uint) error {
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

	nodeConnections[node.ID] = conn
}

func (ns *Node) GetNodeConnection(id uint) *websocket.Conn {
	nodeLocker.Lock()
	defer nodeLocker.Unlock()

	return nodeConnections[id]
}
