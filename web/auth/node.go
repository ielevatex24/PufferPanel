package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferpanel/v3/middleware/panelmiddleware"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"net/http"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func registerNodeLogin(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	db := panelmiddleware.GetDatabase(c)

	ns := services.Node{DB: db}

	node, err := ns.Get(1)
	ns.AddNodeConnection(node, conn)
}
