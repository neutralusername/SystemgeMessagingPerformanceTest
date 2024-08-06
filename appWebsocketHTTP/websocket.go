package appWebsocketHTTP

import (
	"SystemgeMessagingPerformanceTest/topics"
	"sync/atomic"
	"time"

	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{
		topics.ASYNC: func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			startedAt := time.Now()
			for i := 0; i < 100000; i++ {
				func() {
					err := node.AsyncMessage(topics.ASYNC, "", "nodeApp")
					if err != nil {
						panic(err)
					}
				}()
			}
			println("100000 async messages sent in " + time.Since(startedAt).String())
			return nil
		},
		topics.SYNC: func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			counter := atomic.Uint32{}
			startedAt := time.Now()
			for i := 0; i < 100000; i++ {
				func() {
					if responseChannel, err := node.SyncMessage(topics.SYNC, "", "nodeApp"); err != nil {
						panic(err)
					} else {
						go func(responseChannel *Node.SyncResponseChannel) {
							if _, err = responseChannel.ReceiveResponse(); err != nil {
								panic(err)
							}
							counter := counter.Add(1)
							if counter == 100000 {
								println("100000 sync responses received in " + time.Since(startedAt).String())
							}
						}(responseChannel)
					}
				}()
			}
			println("100000 sync requests sent in " + time.Since(startedAt).String())
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {

}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
}
