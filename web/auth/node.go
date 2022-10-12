package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferpanel/v3/middleware/panelmiddleware"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func registerNodeLogin(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Node ") {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	db := panelmiddleware.GetDatabase(c)
	ss := &services.Session{DB: db}
	ns := services.Node{DB: db}

	code := strings.TrimPrefix(authHeader, "Node ")

	node, err := ss.ValidateNode(code)

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	} else if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		//upgrader already replies to the client, so we just abort
		c.Abort()
		return
	}

	ns.AddNodeConnection(node, conn)
}
