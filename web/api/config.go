package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferpanel/v3/config"
	"net/http"
	"os"
	"strings"
)

func panelConfig(c *gin.Context) {
	themes := []string{}
	files, err := os.ReadDir(config.WebRoot.Value() + "/theme")
	if err != nil {
		themes = append(themes, "PufferPanel")
	} else {
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".tar") {
				themes = append(themes, f.Name()[:len(f.Name())-4])
			}
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"themes": map[string]interface{}{
			"active":    config.DefaultTheme.Value(),
			"settings":  config.ThemeSettings.Value(),
			"available": themes,
		},
		"branding": map[string]interface{}{
			"name": config.CompanyName.Value(),
		},
		"registrationEnabled": config.RegistrationEnabled.Value(),
	})
}
