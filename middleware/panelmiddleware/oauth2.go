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

package panelmiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"net/http"
)

func HasPermission(requiredScope pufferpanel.Scope, requireServer bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := GetDatabase(c)

		if db == nil {
			NeedsDatabase(c)
			db = GetDatabase(c)
			if db == nil {
				response.HandleError(c, pufferpanel.ErrDatabaseNotAvailable, http.StatusInternalServerError)
				return
			}
		}

		ss := &services.Server{DB: db}
		ps := &services.Permission{DB: db}

		u, exists := c.Get("user")
		if !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, ok := u.(*models.User)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var serverId string

		var server *models.Server
		var err error

		i := c.Param("serverId")
		if requireServer {
			server, err = ss.Get(i)
			if err != nil {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		if requireServer && (server == nil || server.Identifier == "") {
			c.AbortWithStatus(http.StatusForbidden)
			return
		} else if requireServer {
			serverId = server.Identifier
		}

		allowed := false

		if requiredScope != pufferpanel.ScopeNone {
			permissions, err := ps.GetForUserAndServer(user.ID, &serverId)
			if err != nil {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			if pufferpanel.ContainsScope(permissions.ToScopes(), requiredScope) {
				allowed = true
			} else if serverId != "" {
				//if there isn't a defined rule, is this user an admin?
				permissions, err = ps.GetForUserAndServer(user.ID, &serverId)
				if err != nil {
					c.AbortWithStatus(http.StatusForbidden)
					return
				}
				if pufferpanel.ContainsScope(permissions.ToScopes(), pufferpanel.ScopeServersAdmin) {
					allowed = true
				}
			}
		} else {
			allowed = true
		}

		if !allowed {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set("server", server)
	}
}
