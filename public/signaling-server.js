// Simple signaling server (Node.js + WebSocket)
const WebSocket = require('ws');
const wss = new WebSocket.Server({ port: 8080 });

let clients = [];

wss.on('connection', (ws) => {
  clients.push(ws);
  console.log('🔌 New client connected');

  ws.on('message', (message) => {
    console.log('📩 Received:', message.toString());

    // Send to all other clients
    clients.forEach(client => {
      if (client !== ws && client.readyState === WebSocket.OPEN) {
        client.send(message);
      }
    });
  });

  ws.on('close', () => {
    console.log('❌ Client disconnected');
    clients = clients.filter(c => c !== ws);
  });
});

console.log('✅ Signaling server running on ws://localhost:8080');

