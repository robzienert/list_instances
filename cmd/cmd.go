// Package cmd defines the common, application-level API for the tool.
package cmd

import "gopkg.in/alecthomas/kingpin.v2"

var (
	// TODO Override by envar.
	defaultRegions = []string{
		"us-east-1",
		"us-west-2",
		"eu-west-1",
		"ap-northeast-1",
		"ap-northeast-2",
	}
)

// ApplicationOptions defines the global application config options.
type ApplicationOptions struct {
	Debug   bool
	Regions []string
	Retries int
}

// CreateApplicationOptions defines the application flags and binds the kingpin
// values to the opts struct.
func CreateApplicationOptions(app *kingpin.Application) *ApplicationOptions {
	o := &ApplicationOptions{}
	app.Flag("debug", "Enable debug mode.").BoolVar(&o.Debug)
	app.Flag("region", "AWS Region to query. Repeat for multiple regions.").Short('r').Default(defaultRegions...).StringsVar(&o.Regions)
	app.Flag("retries", "Number of times to retry AWS API calls in case of errors.").Default("10").IntVar(&o.Retries)
	return o
}
