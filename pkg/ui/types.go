package ui

import "github.com/sh31k30ps/gikopsctl/pkg/config"

type UIRequestResult interface{}

type UIRequester interface {
	Request() (UIRequestResult, error)
}

type UIRequesterConfigAware interface {
	Config() config.ConfigObject
}
