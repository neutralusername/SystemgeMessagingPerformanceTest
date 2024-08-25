package app

import (
	"SystemgeMessagingPerformanceTest/topics"
	"sync/atomic"
	"time"

	"github.com/neutralusername/Systemge/Commands"
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SystemgeClient"
	"github.com/neutralusername/Systemge/SystemgeConnection"
)

var async_startedAt = time.Time{}
var async_counter = atomic.Uint32{}

var sync_startedAt = time.Time{}
var sync_counter = atomic.Uint32{}

type App struct {
	systemgeClient *SystemgeClient.SystemgeClient
}

func New() *App {
	app := &App{}

	messageHandler := SystemgeConnection.NewConcurrentMessageHandler(
		SystemgeConnection.AsyncMessageHandlers{
			topics.ASYNC: func(connection *SystemgeConnection.SystemgeConnection, message *Message.Message) {
				val := sync_counter.Add(1)
				if val == 1 {
					sync_startedAt = time.Now()
				}
				if val == 100000 {
					println("100000 async messages received in " + time.Since(sync_startedAt).String())
					sync_counter.Store(0)
				}
			},
		},
		SystemgeConnection.SyncMessageHandlers{
			topics.SYNC: func(connection *SystemgeConnection.SystemgeConnection, message *Message.Message) (string, error) {
				val := async_counter.Add(1)
				if val == 1 {
					async_startedAt = time.Now()
				}
				if val == 100000 {
					println("100000 sync requests received in " + time.Since(async_startedAt).String())
					async_counter.Store(0)
				}
				return "", nil
			},
		},
		nil, nil,
	)
	app.systemgeClient = SystemgeClient.New(
		&Config.SystemgeClient{
			Name: "systemgeClient",
			EndpointConfigs: []*Config.TcpEndpoint{
				{
					Address: "localhost:60001",
				},
			},
			ConnectionConfig: &Config.SystemgeConnection{},
		},
		func(connection *SystemgeConnection.SystemgeConnection) error {
			connection.StartProcessingLoopSequentially(messageHandler)
			return nil
		},
		func(connection *SystemgeConnection.SystemgeConnection) {
			connection.StopProcessingLoop()
		},
	)
	Dashboard.NewClient(
		&Config.DashboardClient{
			Name:             "appGameOfLife",
			ConnectionConfig: &Config.SystemgeConnection{},
			EndpointConfig: &Config.TcpEndpoint{
				Address: "localhost:60000",
			},
		},
		app.systemgeClient.Start, app.systemgeClient.Stop, app.systemgeClient.GetMetrics, app.systemgeClient.GetStatus,
		Commands.Handlers{})
	return app
}
