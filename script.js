const ws = new WebSocket("ws://localhost:8080/ws");

const messagesDiv = document.getElementById("messages");
const usernameInput = document.getElementById("username");
const messageInput = document.getElementById("message");
const sendBtn = document.getElementById("sendBtn");

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  const div = document.createElement("div");
  div.textContent = `${msg.username}: ${msg.content}`;
  messagesDiv.appendChild(div);
};

sendBtn.onclick = () => {
  const msg = {
    username: usernameInput.value || "Người dùng",
    content: messageInput.value
  };
  ws.send(JSON.stringify(msg));
  messageInput.value = "";
};
