package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mp4096/bart"
)

const (
	ok                      int = 0
	noConfigFileSpecified   int = 3
	noTemplateFileSpecified int = 4
	couldNotOpenConfigFile  int = 5
	couldNotSendOrPreview   int = 6
)

var send = false
var fs = flag.NewFlagSet("bart", flag.ExitOnError)
var configFilename = ""
var templateFilename = ""
var revision = "undefined revision"
var buildTime = "undefined build time"

func init() {
	fs.BoolVar(&send, "s", false, "Send emails; dry run otherwise")
	fs.StringVar(&configFilename, "c", "", "Config filename")
	fs.StringVar(&templateFilename, "t", "", "Template filename")
}

func main() {
	if len(os.Args) == 1 {
		printHelp()
		os.Exit(ok)
	}

	fs.Parse(os.Args[1:])

	if len(configFilename) == 0 {
		fmt.Println("Config file not specified")
		os.Exit(noConfigFileSpecified)
	}

	if len(templateFilename) == 0 {
		fmt.Println("Template file not specified")
		os.Exit(noTemplateFileSpecified)
	}

	c := new(bart.Config)
	errConf := c.ImportFromFile(configFilename)
	if errConf != nil {
		fmt.Println("Error opening config file")
		fmt.Println(errConf)
		os.Exit(couldNotOpenConfigFile)
	}

	fmt.Println("Hello,", c.Author.Name)

	if err := bart.ProcessFile(templateFilename, send, c); err != nil {
		fmt.Println("Error while sending email or opening preview")
		fmt.Println(err)
		os.Exit(couldNotSendOrPreview)
	}
}

func printHelp() {
	fmt.Println("bart (" + buildTime + " " + revision + ")")
	fmt.Println("bart sends templated emails for you")
	fmt.Println("Run bart -h for help")
}
