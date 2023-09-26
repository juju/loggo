package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/juju/loggo"
	"github.com/juju/loggo/attrs"
	loggoslog "github.com/juju/loggo/slog"
)

var rootLogger = loggo.GetLogger("")

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println("Add a parameter to configure the logging:")
		fmt.Println(`E.g. "<root>=INFO;first=TRACE" or "<root>=INFO;first=TRACE" "slog"`)
	}
	num := len(args)
	if num > 1 {
		if err := loggo.ConfigureLoggers(args[1]); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("\nCurrent logging levels:")
	fmt.Println(loggo.LoggerInfo())

	if num > 2 {
		if args[2] == "slog" {
			handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: loggoslog.DefaultLevel(loggo.DefaultContext().Config()),
			})
			loggo.ReplaceDefaultWriter(loggoslog.NewSlogWriter(handler))

			fmt.Println("Using log/slog writer:")
		} else {
			log.Fatalf("unknown logging type %q", args[2])
		}
	}

	fmt.Println("")

	rootLogger.Info("Start of test.", attrs.String("foo", "bar"))

	FirstCritical("first critical")
	FirstError("first error")
	FirstWarning("first warning")
	FirstInfo("first info")
	FirstDebug("first debug")
	FirstTrace("first trace")

	SecondCritical("second critical")
	SecondError("second error")
	SecondWarning("second warning")
	SecondInfo("second info")
	SecondDebug("second debug")
	SecondTrace("second trace")
}
