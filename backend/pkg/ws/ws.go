package serversk

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Accepting all requests
	},
}

type Server struct {
	Ver     *sync.Mutex
	clients map[*websocket.Conn]bool
	Channel chan []byte
}

func StartServer() *Server {
	server := Server{
		&sync.Mutex{},
		make(map[*websocket.Conn]bool),
		make(chan []byte),
	}

	http.HandleFunc("/", server.echo)
	go http.ListenAndServe(":8080", nil)

	return &server
}

func (server *Server) echo(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)

	server.clients[connection] = true // Save the connection using it as a key

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break // Exit the loop if the client tries to close the connection or the connection is interrupted
		}
		server.Channel <- message // Send the message to the channel
	}

	delete(server.clients, connection) // Removing the connection

	connection.Close()
}

func (server *Server) WriteMessage(message []byte) {
	for conn := range server.clients {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
