package main

import (
	"fmt"
	"os"

	"github.com/sh31k30ps/gikopsctl/pkg/cmd"
	"github.com/sh31k30ps/gikopsctl/pkg/cmd/gikopsctl"
)

func main() {
	rootCmd := gikopsctl.NewRootCmd(cmd.NewLogger())
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
