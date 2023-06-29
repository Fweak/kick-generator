package kick

import (
	"encoding/json"
	"log"

	"net/http"

	"github.com/gorilla/websocket"
)

type socketReply struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func connectToWebsocket() *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial("wss://ws-us2.pusher.com/app/eb1d5f283081a78b932c?protocol=7&client=js&version=7.6.0&flash=false", http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/113.0"},
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return conn
}

// StartSocket - connects to our sokcet
func (client *Client) StartSocket() {
	client.Conn = connectToWebsocket()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		var payload socketReply
		err = json.Unmarshal(message, &payload)
		if err != nil {
			log.Println(err)
			return
		}

		if payload.Event == "pusher:connection_established" {
			var data map[string]interface{}
			json.Unmarshal([]byte(payload.Data), &data)
			client.socketID = data["socket_id"].(string)
			break
		}
		continue
	}
}
