package server

import (
	"math/rand/v2"
	"net/http"
	"time"
)

func (server *Server) HardOpHandler(response http.ResponseWriter, _ *http.Request) {
	operationTime := rand.IntN(11) + 10
	time.Sleep(time.Duration(operationTime) * time.Second)

	if rand.IntN(2) == 1 {
		response.WriteHeader(http.StatusOK)
		return
	}

	serverErrorCodes := []int{
		http.StatusInternalServerError,
		http.StatusServiceUnavailable,
	}
	response.WriteHeader(serverErrorCodes[rand.IntN(len(serverErrorCodes))])
}
