/*
 Copyright 2020 Padduck, LLC
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

package proxy

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/middleware/panelmiddleware"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

var serverProxy = func(c *gin.Context) { proxyRequest(c, true) }
var nodeProxy = func(c *gin.Context) { proxyRequest(c, false) }

func RegisterRoutes(rg *gin.RouterGroup) {
	proxy := rg.Group("/daemon", panelmiddleware.AuthMiddleware, panelmiddleware.NeedsDatabase)
	{
		g := proxy.Group("/server")
		{
			g.Any("/:id", serverProxy)
			g.Any("/:id/*path", serverProxy)
		}

		g = proxy.Group("/socket")
		{
			g.Any("/:id", serverProxy)
		}
	}

	proxy = rg.Group("/node", panelmiddleware.AuthMiddleware, panelmiddleware.NeedsDatabase)
	{
		proxy.Any("/:id", nodeProxy)
		proxy.Any("/:id/*path", nodeProxy)
	}
}

func proxyRequest(c *gin.Context, isServer bool) {
	db := panelmiddleware.GetDatabase(c)
	ns := &services.Node{DB: db}

	path := strings.TrimPrefix(c.Request.URL.Path, "/proxy")
	var node *models.Node

	if isServer {
		ss := &services.Server{DB: db}

		serverId := c.Param("id")
		if serverId == "" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		s, err := ss.Get(serverId)
		if err != nil && gorm.ErrRecordNotFound != err && response.HandleError(c, err, http.StatusInternalServerError) {
			return
		} else if s == nil || s.Identifier == "" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		node = &s.Node
	} else {
		nodeId := c.Param("id")
		//remove the node's id from the path
		path = "/daemon" + strings.TrimPrefix(path, "/node/"+nodeId)

		id, err := cast.ToUintE(nodeId)
		if response.HandleError(c, err, http.StatusBadRequest) {
			return
		}

		node, err = ns.Get(id)
		if response.HandleError(c, err, http.StatusBadRequest) {
			return
		}
	}

	if node.IsLocal() {
		c.Request.URL.Path = path
		pufferpanel.Engine.HandleContext(c)
	} else {
		if c.IsWebsocket() {
			proxySocketRequest(c, path, ns, node)
		} else {
			proxyHttpRequest(c, path, ns, node)
		}
	}

	c.Abort()
}

func proxyHttpRequest(c *gin.Context, path string, ns *services.Node, node *models.Node) {
	callResponse, err := ns.CallNode(node, c.Request.Method, path, c.Request.Body, c.Request.Header)

	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	//Even though apache isn't going to be in place, we can't set certain headers
	newHeaders := make(map[string]string, 0)
	for k, v := range callResponse.Header {
		switch k {
		case "Transfer-Encoding":
		case "Content-Type":
		case "Content-Length":
			continue
		default:
			newHeaders[k] = strings.Join(v, ", ")
		}
	}

	c.DataFromReader(callResponse.StatusCode, callResponse.ContentLength, callResponse.Header.Get("Content-Type"), callResponse.Body, newHeaders)
	c.Abort()
}

func proxySocketRequest(c *gin.Context, path string, ns *services.Node, node *models.Node) {
	if node.IsLocal() {
		//have gin handle the request again, but send it to daemon instead
		pufferpanel.Engine.HandleContext(c)
	} else {
		err := ns.OpenSocket(node, path, c.Writer, c.Request)
		response.HandleError(c, err, http.StatusInternalServerError)
	}
	c.Abort()
}
