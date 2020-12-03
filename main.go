package main

import (
	"os"

	"github.com/sam-lord/varnishbeat/cmd"

	// Make sure all your modules and metricsets are linked in this file
	_ "github.com/sam-lord/varnishbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
