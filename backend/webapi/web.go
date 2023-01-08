package webapi

import (
	"encoding/json"
	"fmt"
	"github.com/fafeitsch/city-knowledge-contest/backend/contest"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

var openRooms = make(map[string]*contest.Room)

func parseMessage[K any](message json.RawMessage) K {
	var request K
	_ = json.Unmarshal(message, &request)
	return request
}

func HandleFunc(options Options) http.HandlerFunc {
	methods := map[string]rpcHandler{
		"createRoom":     createRoom,
		"updateRoom":     updateRoom,
		"joinRoom":       joinRoom,
		"startGame":      startGame,
		"answerQuestion": answerQuestion,
		"advanceGame":    advanceGame,
	}
	return func(resp http.ResponseWriter, req *http.Request) {
		if options.AllowCors {
			setCorsHeaders(resp)
		}
		if req.Method == "OPTIONS" {
			return
		}
		if req.Header.Get("Upgrade") == "websocket" {
			err := handleWebsocketRequest(resp, req, options)
			if err != nil {
				fmt.Printf("%v", err)
				resp.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		if req.Method != "POST" {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		parts := strings.Split(req.RequestURI, "/")
		if len(parts) != 2 || parts[1] != "rpc" {
			resp.WriteHeader(http.StatusNotFound)
			return
		}
		var request Request
		defer handlePanic(resp, request)
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			writeError(resp, -32700, nil, "could not parse JSON-RPC request: %v", err)
			return
		}
		validator, ok := methods[request.Method]
		if !ok {
			writeError(resp, -32601, request.Id, "the requested method \"%s\" was not found", request.Method)
			return
		}
		rpcRequest, err := validator(request.Params)
		if rpcRequest != nil && rpcRequest.release != nil {
			defer rpcRequest.release()
		}
		if err != nil {
			writeError(resp, -32602, request.Id, "the validation for method \"%s\" failed: %v", request.Method, err)
			return
		}
		result, err := rpcRequest.process()
		if err != nil {
			writeError(resp, -32603, request.Id, "execution of method \"%s\" failed: %v", request.Method, err)
			return
		}
		jsonResult, _ := json.Marshal(result)
		response := Response{Id: request.Id, Jsonrpc: "2.0", Result: jsonResult}
		_ = json.NewEncoder(resp).Encode(response)
	}
}

func setCorsHeaders(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func handlePanic(resp http.ResponseWriter, request Request) {
	if r := recover(); r != nil {
		writeError(resp, -32603, request.Id, "a fatal error occurred during the method execution.")
		log.Printf("panic during RPC call: %v\n%v", r, string(debug.Stack()))
	}
}

func writeError(resp http.ResponseWriter, code int, id *string, format string, params ...interface{}) {
	msg := fmt.Sprintf(format, params...)
	response := Response{
		Id:      id,
		Error:   &Error{Code: code, Message: msg},
		Jsonrpc: "2.0",
	}
	_ = json.NewEncoder(resp).Encode(response)
}
