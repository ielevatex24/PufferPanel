package extract

import (
	"github.com/pufferpanel/pufferpanel/v3/daemon"
)

type OperationFactory struct {
	daemon.OperationFactory
}

func (of OperationFactory) Key() string {
	return "extract"
}
func (of OperationFactory) Create(op daemon.CreateOperation) (daemon.Operation, error) {
	source := op.OperationArgs["source"].(string)
	destination := op.OperationArgs["destination"].(string)

	return Extract{
		Source:      source,
		Destination: destination,
	}, nil
}

var Factory OperationFactory
