package server

import "net/http"

func (server *Server) VersionHandler(response http.ResponseWriter, _ *http.Request) {
	type VersionResponse struct {
		Version string `json:"version"`
	}
	server.Answer(response, &VersionResponse{"v1.0.0"})
}
