package main

import (
	"SystemgeMessagingPerformanceTest/app"
	"SystemgeMessagingPerformanceTest/appWebsocketHttp"
	"time"
)

const LOGGER_PATH = "logs.log"

func main() {
	appWebsocketHttp.New()
	app.New()
	<-make(chan time.Time)
}
