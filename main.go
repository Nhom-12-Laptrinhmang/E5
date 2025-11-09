package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

// Message định nghĩa cấu trúc message
type Message struct {
    Username string `json:"username"`
    Content  string `json:"content"`
}

var (
    clients   = make(map[*websocket.Conn]bool)
    broadcast = make(chan Message)
    upgrader  = websocket.Upgrader{}
    mu        sync.Mutex
)

func main() {
    fs := http.FileServer(http.Dir("./static")) // thư mục chứa index.html, JS, CSS
    http.Handle("/", fs)

    http.HandleFunc("/ws", handleConnections)

    go handleMessages()

    log.Println("Server bắt đầu trên :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}

// Xử lý kết nối WebSocket
func handleConnections(w http.ResponseWriter, r *http.Request) {
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Lỗi upgrade: %v", err)
        return
    }
    defer ws.Close()

    mu.Lock()
    clients[ws] = true
    mu.Unlock()

    log.Printf("Client mới kết nối. Tổng client: %d", len(clients))

    for {
        mt, p, err := ws.ReadMessage()
        if err != nil {
            if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
                log.Printf("Client đóng kết nối bình thường: %v", err)
            } else {
                log.Printf("Lỗi đọc message: %v", err)

            }
            mu.Lock()
            delete(clients, ws)
            mu.Unlock()
            break
        }

        var msg Message
        if mt == websocket.TextMessage {
            if err := json.Unmarshal(p, &msg); err != nil {
                // Nếu không phải JSON -> coi payload là plain text
                msg = Message{
                    Username: "Anon",
                    Content:  string(p),
                }
            }
        } else {
            msg = Message{
                Username: "Anon",
                Content:  fmt.Sprintf("[binary %d bytes]", len(p)),
            }
        }

        broadcast <- msg
    }
}

// Gửi message tới tất cả client
func handleMessages() {
    for {
        msg := <-broadcast
        mu.Lock()
        for client := range clients {
            err := client.WriteJSON(msg)
            if err != nil {
                log.Printf("Lỗi gửi message tới client: %v", err)
                client.Close()
                delete(clients, client)
            }
        }
        mu.Unlock()
    }
}
