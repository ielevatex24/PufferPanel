package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3/middleware/panelmiddleware"
	"github.com/pufferpanel/pufferpanel/v3/response"
	"github.com/pufferpanel/pufferpanel/v3/services"
	"net/http"
)

func LogoutPost(c *gin.Context) {
	db := panelmiddleware.GetDatabase(c)
	ss := services.Session{DB: db}

	cookie, err := c.Cookie("puffer_auth")
	if response.HandleError(c, err, http.StatusInternalServerError) {
		return
	}

	_ = ss.Expire(cookie)
	c.Status(http.StatusNoContent)
}
