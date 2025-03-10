package cmd

import (
	"io"
	"os"

	"github.com/sh31k30ps/gikopsctl/pkg/internal/cli"
	"github.com/sh31k30ps/gikopsctl/pkg/internal/env"
	"github.com/sh31k30ps/gikopsctl/pkg/log"
)

// NewLogger returns the standard logger used by the kind CLI
// This logger writes to os.Stderr
func NewLogger() log.Logger {
	var writer io.Writer = os.Stderr
	if env.IsSmartTerminal(writer) {
		writer = cli.NewSpinner(writer)
	}
	return cli.NewLogger(writer, 0)
}

// ColorEnabled returns true if color is enabled for the logger
// this should be used to control output
func ColorEnabled(logger log.Logger) bool {
	type maybeColorer interface {
		ColorEnabled() bool
	}
	v, ok := logger.(maybeColorer)
	return ok && v.ColorEnabled()
}
