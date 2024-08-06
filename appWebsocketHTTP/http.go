package appWebsocketHTTP

import (
	"net/http"

	"github.com/neutralusername/Systemge/HTTP"
)

func (app *AppWebsocketHTTP) GetHTTPMessageHandlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/": HTTP.SendDirectory("../frontend"),
	}
}
