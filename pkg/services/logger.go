package services

import (
	"io"
	"os"

	"github.com/sh31k30ps/gikops/pkg/cli"
	"github.com/sh31k30ps/gikops/pkg/internal/env"
	"github.com/sh31k30ps/gikops/pkg/log"
)

var loggers map[string]log.Logger

// NewLogger returns the standard logger used by the kind CLI
// This logger writes to os.Stderr
func NewLogger(name string) log.Logger {
	if loggers == nil {
		loggers = make(map[string]log.Logger)
	}
	if loggers[name] == nil {
		var writer io.Writer = os.Stdout
		if env.IsSmartTerminal(writer) {
			writer = cli.NewSpinner(writer)
		}
		loggers[name] = cli.NewLogger(writer, 0)
	}
	return loggers[name]
}
