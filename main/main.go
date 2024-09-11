package main

import (
	"SystemgeMessagingPerformanceTest/app"
	"SystemgeMessagingPerformanceTest/appWebsocketHttp"
	"time"

	"github.com/neutralusername/Systemge/Helpers"
)

const LOGGER_PATH = "logs.log"

func main() {
	Helpers.StartPprof()
	appWebsocketHttp.New()
	app.New()
	<-make(chan time.Time)
}
