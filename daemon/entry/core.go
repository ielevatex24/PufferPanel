package entry

import (
	"github.com/pufferpanel/pufferpanel/v3/config"
	"github.com/pufferpanel/pufferpanel/v3/daemon/environments"
	"github.com/pufferpanel/pufferpanel/v3/daemon/programs"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"github.com/pufferpanel/pufferpanel/v3/sftp"
	"os"
	"path/filepath"
	"strings"
)

func Start() error {
	sftp.Run()

	environments.LoadModules()
	programs.Initialize()

	var err error

	if _, err = os.Stat(config.ServersFolder.Value()); os.IsNotExist(err) {
		logging.Info.Printf("No server directory found, creating")
		err = os.MkdirAll(config.ServersFolder.Value(), 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}
	}

	//update path to include our binary folder
	newPath := os.Getenv("PATH")
	fullPath, _ := filepath.Abs(config.BinariesFolder.Value())
	if !strings.Contains(newPath, fullPath) {
		_ = os.Setenv("PATH", newPath+":"+fullPath)
	}

	programs.LoadFromFolder()

	programs.InitService()

	for _, element := range programs.GetAll() {
		if element.IsEnabled() {
			element.GetEnvironment().DisplayToConsole(true, "Daemon has been started\n")
			if element.IsAutoStart() {
				logging.Info.Printf("Queued server %s", element.Id())
				element.GetEnvironment().DisplayToConsole(true, "Server has been queued to start\n")
				programs.StartViaService(element)
			}
		}
	}
	return nil
}
