package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/go-acme/lego/v4/cmd"
	"github.com/go-acme/lego/v4/log"
	"github.com/jasonlvhit/gocron"
	"github.com/urfave/cli"
)

const (
	envEveryXDays = "EVERY_X_DAYS"
)

var (
	version    = "dev"
	app        = cli.NewApp()
	everyXDays = 1
)

func init() {
	app.Name = "lego"
	app.HelpName = "lego"
	app.Usage = "Let's Encrypt client written in Go"
	app.EnableBashCompletion = true

	app.Version = version
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("lego version %s %s/%s\n", c.App.Version, runtime.GOOS, runtime.GOARCH)
	}

	var defaultPath string
	cwd, err := os.Getwd()
	if err == nil {
		defaultPath = filepath.Join(cwd, ".lego")
	}

	app.Flags = cmd.CreateFlags(defaultPath)
	app.Before = cmd.Before
	app.Commands = cmd.CreateCommands()

	e := os.Getenv(envEveryXDays)
	if e == "" {
		everyXDays = 1
		log.Println("not set env:", envEveryXDays, ". use default value every 1 day")
	} else {
		everyXDays, err := strconv.Atoi(e)
		if err != nil {
			everyXDays = 1
			log.Println("env:", envEveryXDays, " is not a number, use default value every 1 day")
		} else if everyXDays < 1 {
			everyXDays = 1
			log.Println("env:", envEveryXDays, " is less than or equal to zero, use default value every 1 day")
		}
	}
}

func task() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func main() {
	s := gocron.NewScheduler()
	if everyXDays == 1 {
		s.Every(uint64(everyXDays)).Day().Do(task)
	} else {
		s.Every(uint64(everyXDays)).Days().Do(task)
	}
	<-s.Start()
}
