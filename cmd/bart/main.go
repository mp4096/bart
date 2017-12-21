package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mp4096/bart"
)

const (
	OK                         int = 0
	INVALID_COMMAND            int = 1
	NO_CONFIG_FILE_SPECIFIED   int = 3
	NO_TEMPLATE_FILE_SPECIFIED int = 4
	COULD_NOT_OPEN_CONFIG_FILE int = 5
	COULD_NOT_SEND_OR_PREVIEW  int = 6
)

var send = false
var fs = flag.NewFlagSet("bart", flag.ExitOnError)
var configFilename = ""
var templateFilename = ""

func init() {
	fs.BoolVar(&send, "s", false, "Send emails; dry run otherwise")
	fs.StringVar(&configFilename, "c", "", "Config filename")
	fs.StringVar(&templateFilename, "t", "", "Template filename")
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("bart sends templated emails for you")
		fmt.Println("Run bart -h for help")
		os.Exit(OK)
	}

	fs.Parse(os.Args[1:])

	if len(configFilename) == 0 {
		fmt.Println("Config file not specified")
		os.Exit(NO_CONFIG_FILE_SPECIFIED)
	}

	if len(templateFilename) == 0 {
		fmt.Println("Template file not specified")
		os.Exit(NO_TEMPLATE_FILE_SPECIFIED)
	}

	c := new(bart.Config)
	errConf := c.ImportFromFile(configFilename)
	if errConf != nil {
		fmt.Println("Error opening config file")
		fmt.Println(errConf)
		os.Exit(COULD_NOT_OPEN_CONFIG_FILE)
	}

	fmt.Println("Hello,", c.Author.Name)

	if err := bart.ProcessFile(templateFilename, send, c); err != nil {
		fmt.Println("Error while sending email or opening preview")
		fmt.Println(err)
		os.Exit(COULD_NOT_SEND_OR_PREVIEW)
	}
}
