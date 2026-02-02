package handler

import "github.com/gustavoz65/Cashing-go/internal/server"

type Handler struct {
	server *server.Server
}

func NewHandler(s *server.Server) Handler {
	return Handler{server: s}
}
