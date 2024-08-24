package main

import (
	"SystemgeMessagingPerformanceTest/app"
	"SystemgeMessagingPerformanceTest/appWebsocketHttp"
	"time"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
)

const LOGGER_PATH = "logs.log"

func main() {
	Dashboard.NewServer(&Config.DashboardServer{
		HTTPServerConfig: &Config.HTTPServer{
			TcpListenerConfig: &Config.TcpListener{
				Port: 8081,
			},
		},
		WebsocketServerConfig: &Config.WebsocketServer{
			Pattern:                 "/ws",
			ClientWatchdogTimeoutMs: 1000 * 60,
			TcpListenerConfig: &Config.TcpListener{
				Port: 8444,
			},
		},
		SystemgeServerConfig: &Config.SystemgeServer{
			Name: "dashboardServer",
			ListenerConfig: &Config.SystemgeListener{
				TcpListenerConfig: &Config.TcpListener{
					Port: 60000,
				},
			},
			ConnectionConfig: &Config.SystemgeConnection{},
		},
		HeapUpdateIntervalMs:      1000,
		GoroutineUpdateIntervalMs: 1000,
		StatusUpdateIntervalMs:    1000,
		MetricsUpdateIntervalMs:   1000,
	})
	appWebsocketHttp.New()
	app.New()
	<-make(chan time.Time)
}
