package main

import (
	"fmt"
	"os"

	"github.com/juju/loggo"
)

var logger = loggo.GetLogger("main")
var rootLogger = loggo.GetLogger("")

func main() {
	args := os.Args
	if len(args) > 1 {
		loggo.ConfigureLoggers(args[1])
	} else {
		fmt.Println("Add a parameter to configure the logging:")
		fmt.Println("E.g. \"<root>=INFO;first=TRACE\"")
	}
	fmt.Println("\nCurrent logging levels:")
	fmt.Println(loggo.LoggerInfo())
	fmt.Println("")

	rootLogger.Infof("Start of test.")

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
