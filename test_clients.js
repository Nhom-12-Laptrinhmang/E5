const N = 10; // số lượng client ảo
const clients = [];

// Bước 1: Tạo kết nối WebSocket
for (let i = 0; i < N; i++) {
  const ws = new WebSocket("ws://localhost:8080/ws");
  ws.onopen = () => {
    ws.send(`Client ${i} connected`);
  };
  ws.onmessage = (msg) => {
    // console.log(`Client ${i} nhận: ${msg.data}`);
  };
  ws.onclose = () => {
    console.log(`Client ${i} disconnected`);
  };
  ws.onerror = (err) => {
    console.error(`Client ${i} lỗi:`, err);
  };
  clients.push(ws);
}

// Bước 2: Gửi tin nhắn định kỳ từ các client
clients.forEach((ws, i) => {
  setInterval(() => {
    if (ws.readyState === WebSocket.OPEN) {
      ws.send(`Message từ client ${i}`);
    }
  }, 1000); // mỗi 1 giây gửi 1 tin nhắn
});
