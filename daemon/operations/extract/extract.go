package extract

import (
	"github.com/mholt/archiver"
	"github.com/pufferpanel/pufferpanel/v3/daemon"
)

type Extract struct {
	Source      string
	Destination string
}

func (op Extract) Run(daemon.Environment) error {
	return archiver.Unarchive(op.Source, op.Destination)
}
