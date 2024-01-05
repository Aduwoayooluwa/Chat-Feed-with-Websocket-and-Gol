package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	connections map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		connections: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("Hello, we can begin...", ws.RemoteAddr())

	s.connections[ws] = true


	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Something occured", err)
			delete(s.connections, ws)
			break
		}

		s.broadcast(message)
	}
}

func (s *Server) broadcast(message []byte) {
	for conn := range s.connections {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("Error in communicating", err)
		}
	}
}

func main() {
	server := NewServer()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if err != nil {
			fmt.Println("No vex, could not start convo", err)
			return
		}
		defer conn.Close()

		server.handleWS(conn)
	})

	fmt.Println("App is ready at port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Something don sup!", err)
	}
}
