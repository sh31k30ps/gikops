package component

import (
	"github.com/sh31k30ps/gikopsctl/pkg/log"
)

type Command struct {
	logger log.Logger
}

func NewCommand(logger log.Logger) *Command {
	return &Command{logger: logger}
}

func (c *Command) Create() error {
	return nil
}

func (c *Command) Edit() error {
	return nil
}

func (c *Command) Delete(id interface{}) error {
	return nil
}

func (c *Command) Add() error {
	return nil
}
