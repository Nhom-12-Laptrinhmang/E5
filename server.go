package main

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        conn, _ := upgrader.Upgrade(w, r, nil)
        defer conn.Close()

        for {
            _, msg, err := conn.ReadMessage()
            if err != nil { break }
            log.Println("Received:", string(msg))
            conn.WriteMessage(websocket.TextMessage, []byte("Echo: "+string(msg)))
        }
    })

    log.Println("Server running on :8080")
    http.ListenAndServe(":8080", nil)
}