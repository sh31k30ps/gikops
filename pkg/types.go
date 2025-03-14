package pkg

import (
	"fmt"

	"github.com/sh31k30ps/gikopsctl/pkg/cluster"
	"github.com/sh31k30ps/gikopsctl/pkg/component"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
	"github.com/sh31k30ps/gikopsctl/pkg/project"
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
	Edit() error
	Delete(id interface{}) error
	Add() error
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
