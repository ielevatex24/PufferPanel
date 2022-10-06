package archive

import (
	"github.com/mholt/archiver"
	"github.com/pufferpanel/pufferpanel/v3/daemon"
)

type Archive struct {
	Source      []string
	Destination string
}

func (op Archive) Run(daemon.Environment) error {
	return archiver.Archive(op.Source, op.Destination)
}
