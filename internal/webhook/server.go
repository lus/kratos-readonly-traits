package webhook

import (
	"context"
	"net/http"
)

type Server struct {
	httpServer   *http.Server
	Address      string
	ErrorMessage string
}

func (server *Server) Start() error {
	cnt := &controller{
		ErrorMessage: server.ErrorMessage,
	}
	http.HandleFunc("/", cnt.Endpoint)
	httpServer := &http.Server{
		Addr: server.Address,
	}
	server.httpServer = httpServer
	return httpServer.ListenAndServe()
}

func (server *Server) Stop(ctx context.Context) error {
	if server.httpServer == nil {
		return nil
	}
	if err := server.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	server.httpServer = nil
	return nil
}
