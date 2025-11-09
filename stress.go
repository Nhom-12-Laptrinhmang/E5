package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func main() {
	const (
		N       = 100 // số client ảo
		MSG_PER = 5   // số message mỗi client gửi
		SERVER  = "ws://localhost:8080/ws"
	)

	var wg sync.WaitGroup
	wg.Add(N)

	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			ws, _, err := websocket.DefaultDialer.Dial(SERVER, nil)
			if err != nil {
				log.Printf("Client %d: lỗi kết nối: %v", i, err)
				return
			}
			defer ws.Close()

			// Gửi message kết nối đầu tiên
			msg := Message{Username: fmt.Sprintf("Client-%d", i), Content: "connected"}
			ws.WriteJSON(msg)

			// Gửi MSG_PER message, mỗi message cách nhau 500ms
			for j := 0; j < MSG_PER; j++ {
				msg := Message{
					Username: fmt.Sprintf("Client-%d", i),
					Content:  fmt.Sprintf("Message %d từ Client-%d", j, i),
				}
				err := ws.WriteJSON(msg)
				if err != nil {
					log.Printf("Client %d: lỗi gửi message: %v", i, err)
					return
				}
				time.Sleep(500 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	log.Println("Stress test hoàn tất!")
}
