package appWebsocketHttp

import (
	"SystemgeMessagingPerformanceTest/topics"
	"sync"
	"sync/atomic"
	"time"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/HTTPServer"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Status"
	"github.com/neutralusername/Systemge/SystemgeConnection"
	"github.com/neutralusername/Systemge/SystemgeServer"
	"github.com/neutralusername/Systemge/WebsocketServer"
)

type AppWebsocketHTTP struct {
	status      int
	statusMutex sync.Mutex

	systemgeServer  *SystemgeServer.SystemgeServer
	websocketServer *WebsocketServer.WebsocketServer
	httpServer      *HTTPServer.HTTPServer

	connection SystemgeConnection.SystemgeConnection
}

func New() *AppWebsocketHTTP {
	app := &AppWebsocketHTTP{}

	messageHandler := SystemgeConnection.NewTopicExclusiveMessageHandler(
		SystemgeConnection.AsyncMessageHandlers{},
		SystemgeConnection.SyncMessageHandlers{},
		nil, nil, 100000,
	)
	app.systemgeServer = SystemgeServer.New("systemgeServer",
		&Config.SystemgeServer{
			ListenerConfig: &Config.TcpSystemgeListener{
				TcpServerConfig: &Config.TcpServer{
					Port: 60001,
				},
			},
			ConnectionConfig: &Config.TcpSystemgeConnection{},
		},
		func(connection SystemgeConnection.SystemgeConnection) error {
			connection.StartProcessingLoopSequentially(messageHandler)
			app.connection = connection
			return nil
		},
		func(connection SystemgeConnection.SystemgeConnection) {
			connection.StopProcessingLoop()
			app.connection = nil
		},
	)
	app.websocketServer = WebsocketServer.New("websocketServer",
		&Config.WebsocketServer{
			ClientWatchdogTimeoutMs: 1000 * 60,
			Pattern:                 "/ws",
			TcpServerConfig: &Config.TcpServer{
				Port: 8443,
			},
		},

		WebsocketServer.MessageHandlers{
			topics.ASYNC: func(websocketClient *WebsocketServer.WebsocketClient, message *Message.Message) error {
				startedAt := time.Now()
				for i := 0; i < 100000; i++ {
					func() {
						app.connection.AsyncMessage(topics.ASYNC, Helpers.IntToString(i))
					}()
				}
				println("100000 async messages sent in " + time.Since(startedAt).String())
				return nil
			},
			topics.SYNC: func(websocketClient *WebsocketServer.WebsocketClient, message *Message.Message) error {
				counter := atomic.Uint32{}
				startedAt := time.Now()
				for i := 0; i < 100000; i++ {
					func() {
						if responseChannel, err := app.connection.SyncRequest(topics.SYNC, ""); err != nil {
							panic(err)
						} else {
							go func(responseChannel <-chan *Message.Message) {
								response := <-responseChannel
								if response == nil {
									panic("response is nil")
								}
								counter.Add(1)
								if counter.Load() == 100000 {
									println("100000 sync responses received in " + time.Since(startedAt).String())
								}
							}(responseChannel)
						}
					}()
				}
				println("100000 sync requests sent in " + time.Since(startedAt).String())
				return nil
			},
		},
		nil, nil,
	)
	app.httpServer = HTTPServer.New("httpServer",
		&Config.HTTPServer{
			TcpServerConfig: &Config.TcpServer{
				Port: 8080,
			},
		},
		HTTPServer.Handlers{
			"/": HTTPServer.SendDirectory("../frontend"),
		},
	)
	if err := Dashboard.NewClient("appWebsocketHttp",
		&Config.DashboardClient{
			ConnectionConfig: &Config.TcpSystemgeConnection{},
			ClientConfig: &Config.TcpClient{
				Address: "localhost:60000",
			},
		}, app.start, app.stop, app.systemgeServer.GetMetrics, app.getStatus,
		nil,
	).Start(); err != nil {
		panic(err)
	}
	return app
}

func (app *AppWebsocketHTTP) getStatus() int {
	return app.status
}

func (app *AppWebsocketHTTP) start() error {
	app.statusMutex.Lock()
	defer app.statusMutex.Unlock()
	if app.status != Status.STOPPED {
		return Error.New("App already started", nil)
	}
	if err := app.systemgeServer.Start(); err != nil {
		return Error.New("Failed to start systemgeServer", err)
	}
	if err := app.websocketServer.Start(); err != nil {
		app.systemgeServer.Stop()
		return Error.New("Failed to start websocketServer", err)
	}
	if err := app.httpServer.Start(); err != nil {
		app.systemgeServer.Stop()
		app.websocketServer.Stop()
		return Error.New("Failed to start httpServer", err)
	}
	app.status = Status.STARTED
	return nil
}

func (app *AppWebsocketHTTP) stop() error {
	app.statusMutex.Lock()
	defer app.statusMutex.Unlock()
	if app.status != Status.STARTED {
		return Error.New("App not started", nil)
	}
	app.httpServer.Stop()
	app.websocketServer.Stop()
	app.systemgeServer.Stop()
	app.status = Status.STOPPED
	return nil
}
