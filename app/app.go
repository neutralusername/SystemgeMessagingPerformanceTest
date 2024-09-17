package app

import (
	"SystemgeMessagingPerformanceTest/topics"
	"sync/atomic"
	"time"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SystemgeClient"
	"github.com/neutralusername/Systemge/SystemgeConnection"
	"github.com/neutralusername/Systemge/SystemgeMessageHandler"
)

var async_startedAt = time.Time{}
var async_counter = atomic.Uint32{}

var sync_startedAt = time.Time{}
var sync_counter = atomic.Uint32{}

type App struct {
	systemgeClient *SystemgeClient.SystemgeClient
	stopChannel    chan<- bool
}

func New() *App {
	app := &App{}

	messageHandler := SystemgeMessageHandler.NewTopicExclusiveMessageHandler(
		SystemgeMessageHandler.AsyncMessageHandlers{
			topics.ASYNC: func(connection SystemgeConnection.SystemgeConnection, message *Message.Message) {
				val := sync_counter.Add(1)
				if val == 1 {
					sync_startedAt = time.Now()
				}
				if val == 1000000 {
					println("1000000 async messages received in " + time.Since(sync_startedAt).String())
					sync_counter.Store(0)
				}
			},
		},
		SystemgeMessageHandler.SyncMessageHandlers{
			topics.SYNC: func(connection SystemgeConnection.SystemgeConnection, message *Message.Message) (string, error) {
				val := async_counter.Add(1)
				if val == 1 {
					async_startedAt = time.Now()
				}
				if val == 1000000 {
					println("1000000 sync requests received in " + time.Since(async_startedAt).String())
					async_counter.Store(0)
				}
				return "", nil
			},
		},
		nil, nil, 1000000,
	)
	app.systemgeClient = SystemgeClient.New("systemgeClient",
		&Config.SystemgeClient{
			TcpClientConfigs: []*Config.TcpClient{
				{
					Address: "localhost:60001",
				},
			},
			TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{},
		},
		func(connection SystemgeConnection.SystemgeConnection) error {
			stopChannel, _ := SystemgeMessageHandler.StartMessageHandlingLoop_Sequentially(connection, messageHandler)
			app.stopChannel = stopChannel
			return nil
		},
		func(connection SystemgeConnection.SystemgeConnection) {
			close(app.stopChannel)
		},
	)
	if app.systemgeClient.Start() != nil {
		panic("Failed to start systemgeClient")
	}
	return app
}
