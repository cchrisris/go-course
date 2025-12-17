package server

import (
	"encoding/json"
	"net/http"
)

func (server *Server) Answer(response http.ResponseWriter, answerStruct interface{}) {
	server.AnswerWithCode(response, answerStruct, 200)
}

func (server *Server) AnswerWithCode(response http.ResponseWriter,
	answerStruct interface{}, httpCode int) {
	jsonAnswer, err := json.Marshal(answerStruct)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(httpCode)
	_, err = response.Write(jsonAnswer)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}
