package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-acme/lego/v4/cmd"
	"github.com/go-acme/lego/v4/log"
	"github.com/jasonlvhit/gocron"
	"github.com/urfave/cli"
)

const (
	envEveryXDays = "EVERY_X_DAYS"
	envAtTime     = "AT_TIME"
)

var (
	version    = "dev"
	app        = cli.NewApp()
	everyXDays = 1
	atTime     = "01:02:03"
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
	atTime = os.Getenv(envAtTime)
	if e == "" {
		atTime = GetRandomTime()
		log.Println("not set env:", envAtTime, ". use a random time:", atTime)
	} else {
		err = CanFormatTime(atTime)
		if err != nil {
			atTime = GetRandomTime()
			log.Println("not set env:", envAtTime, ". use a random time:", atTime)
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

func CanFormatTime(t string) (err error) {
	var hour, min, sec int
	ts := strings.Split(t, ":")
	if len(ts) < 2 || len(ts) > 3 {
		return err
	}

	if hour, err = strconv.Atoi(ts[0]); err != nil {
		return err
	}
	if min, err = strconv.Atoi(ts[1]); err != nil {
		return err
	}
	if len(ts) == 3 {
		if sec, err = strconv.Atoi(ts[2]); err != nil {
			return err
		}
	}

	if hour < 0 || hour > 23 || min < 0 || min > 59 || sec < 0 || sec > 59 {
		return err
	}

	return nil
}

func GetRandomTime() string {
	var h, m, s int
	rand.Seed(time.Now().UnixNano())
	h = rand.Intn(24)
	m = rand.Intn(60)
	s = rand.Intn(60)
	return strconv.Itoa(h) + ":" + strconv.Itoa(m) + ":" + strconv.Itoa(s)
}

func main() {
	s := gocron.NewScheduler()
	if everyXDays == 1 {
		s.Every(uint64(everyXDays)).Day().At(atTime).Do(task)
	} else {
		s.Every(uint64(everyXDays)).Days().At(atTime).Do(task)
	}
	<-s.Start()
}
