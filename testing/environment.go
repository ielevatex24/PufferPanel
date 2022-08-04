package test

import (
	"errors"
	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/logging"
	"time"
)

type Environment struct {
	*pufferpanel.BaseEnvironment
}

func (*Environment) DisplayToConsole(prefix bool, msg string, data ...interface{}) {
	if prefix {
		logging.Info.Printf("[DAEMON] "+msg, data...)
	} else {
		logging.Info.Printf(msg, data...)
	}
}

func (*Environment) ExecuteInMainProcess(cmd string) (err error) {
	return errors.New("not supported")
}
func (*Environment) Kill() (err error) {
	return errors.New("not supported")
}

func (*Environment) IsRunning() (isRunning bool, err error) {
	return false, errors.New("not supported")
}

func (*Environment) GetStats() (*pufferpanel.ServerStats, error) {
	return nil, errors.New("not supported")
}

func (*Environment) Create() error {
	return errors.New("not supported")
}

func (e *Environment) WaitForMainProcess() error {
	return e.WaitForMainProcessFor(0)
}

func (*Environment) WaitForMainProcessFor(timeout time.Duration) (err error) {
	return errors.New("not supported")
}

func (*Environment) SendCode(code int) error {
	return errors.New("not supported")
}
