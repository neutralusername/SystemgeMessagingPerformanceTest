package app

import (
	"SystemgeMessagingPerformanceTest/topics"
	"sync/atomic"
	"time"

	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

var async_startedAt = time.Time{}
var async_counter = atomic.Uint32{}

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		topics.SYNC: func(node *Node.Node, message *Message.Message) (string, error) {
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
	}
}

var sync_startedAt = time.Time{}
var sync_counter = atomic.Uint32{}

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topics.ASYNC: func(node *Node.Node, message *Message.Message) error {
			val := sync_counter.Add(1)
			if val == 1 {
				sync_startedAt = time.Now()
			}
			if val == 100000 {
				println("100000 async messages received in " + time.Since(sync_startedAt).String())
				sync_counter.Store(0)
			}
			println(message.GetPayload())
			return nil
		},
	}
}
