/*
 Copyright 2019 Padduck, LLC
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

package models

import (
	"github.com/pufferpanel/pufferpanel/v3"
	"gopkg.in/go-playground/validator.v9"
	"net/url"
)

type NodeView struct {
	Id         uint   `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	PublicHost string `json:"publicHost,omitempty"`
	SFTPPort   uint16 `json:"sftpPort,omitempty"`
}

type NodesView []*NodeView

func FromNode(n *Node) *NodeView {
	return &NodeView{
		Id:         n.ID,
		Name:       n.Name,
		PublicHost: n.PublicHost,
		SFTPPort:   n.SFTPPort,
	}
}

func FromNodes(n []*Node) *NodesView {
	result := make(NodesView, len(n))

	for k, v := range n {
		result[k] = FromNode(v)
	}

	return &result
}

func (n *NodeView) CopyToModel(newModel *Node) {
	if n.Name != "" {
		newModel.Name = n.Name
	}

	if n.PublicHost != "" {
		newModel.PublicHost = n.PublicHost
	}

	if n.SFTPPort > 0 {
		newModel.SFTPPort = n.SFTPPort
	}
}

func (n *NodeView) Valid(allowEmpty bool) error {
	validate := validator.New()

	if !allowEmpty && validate.Var(n.Name, "required") != nil {
		return pufferpanel.ErrFieldRequired("name")
	}

	if validate.Var(n.Name, "omitempty,printascii") != nil {
		return pufferpanel.ErrFieldMustBePrintable("name")
	}

	testName := url.QueryEscape(n.Name)
	if testName != n.Name {
		return pufferpanel.ErrFieldHasURICharacters("name")
	}

	if !allowEmpty && validate.Var(n.PublicHost, "required") != nil {
		return pufferpanel.ErrFieldMustBePrintable("publicHost")
	}

	if validate.Var(n.PublicHost, "omitempty,ip|fqdn") != nil {
		return pufferpanel.ErrFieldIsInvalidHost("publicHost")
	}

	if allowEmpty {
		if validate.Var(n.SFTPPort, "min=0,max=65535") != nil {
			return pufferpanel.ErrFieldTooLarge("sftpPort", 65535)
		}
	} else {
		if validate.Var(n.SFTPPort, "min=1,max=65535") != nil {
			return pufferpanel.ErrFieldNotBetween("sftpPort", 1, 65535)
		}
	}

	return nil
}
