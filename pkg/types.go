package pkg

import (
	"fmt"

	"github.com/sh31k30ps/gikops/pkg/cluster"
	"github.com/sh31k30ps/gikops/pkg/component"
	"github.com/sh31k30ps/gikops/pkg/log"
	"github.com/sh31k30ps/gikops/pkg/project"
)

type CommandType string

func (c CommandType) String() string {
	return string(c)
}
func (c CommandType) IsValid() bool {
	for _, cmd := range CommandsTypes {
		if cmd == c {
			return true
		}
	}
	return false
}

const (
	CommandProject   CommandType = "project"
	CommandCluster   CommandType = "cluster"
	CommandComponent CommandType = "component"
)

var (
	CommandsTypes = []CommandType{
		CommandProject,
		CommandCluster,
		CommandComponent,
	}
	CommandsTypesLabels = []string{
		CommandProject.String(),
		CommandCluster.String(),
		CommandComponent.String(),
	}
)

type Command interface {
	Create(args ...interface{}) error
	Edit(mode string, args ...interface{}) error
	Delete(id interface{}) error
	Add(mode string, args ...string) error
	Install() error
}

func GetCommand(Cmd CommandType, logger log.Logger) (Command, error) {
	switch Cmd {
	case CommandProject:
		return project.NewCommand(logger), nil
	case CommandCluster:
		return cluster.NewCommand(logger), nil
	case CommandComponent:
		return component.NewCommand(logger), nil
	default:
		return nil, fmt.Errorf("invalid command: %s", Cmd)
	}
}
