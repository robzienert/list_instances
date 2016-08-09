package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/Sirupsen/logrus"
	"github.com/zerocontribution/list_instances/cmd"
	"github.com/zerocontribution/list_instances/cmd/instances"
)

var (
	app        = kingpin.New("li", "Utility app for additional AWS CLI functionality")
	appOptions = cmd.CreateApplicationOptions(app)
)

func handleDebugFlag(c *kingpin.ParseContext) error {
	if appOptions.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	return nil
}

func main() {
	app.PreAction(handleDebugFlag)
	instances.Bind(app, appOptions)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
