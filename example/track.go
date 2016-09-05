package main

import (
	"github.com/juju/loggo"
	"time"
)

var track = loggo.GetLogger("first")

func FirstTrack(message string) {
	defer track.TimeTrack(loggo.DEBUG, time.Now(), "trackme")
	time.Sleep(time.Second * 2)
	track.Tracef(message)
}
