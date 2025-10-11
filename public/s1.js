const WebSocket = require('ws');
const { v4: uuidv4 } = require('uuid');

// Listen on all interfaces (LAN accessible)
const wss = new WebSocket.Server({ port: 8001, host: 'localhost' });
console.log('✅ WebSocket server running at ws://<your-server-ip>:8001');

const clients = new Map(); // clientId => ws

wss.on('connection', (ws) => {
  const clientId = uuidv4();
  clients.set(clientId, ws);

  console.log(`🔗 Client connected: ${clientId}`);

  // Notify the new client of their ID (optional)
  ws.send(JSON.stringify({ type: "welcome", userId: clientId }));

  // Notify others of new connection (optional)
  broadcastExcept(clientId, {
    type: "user-joined",
    userId: clientId
  });

  // Handle incoming messages and broadcast to others
  ws.on('message', (data) => {
    for (const [id, client] of clients) {
      if (client !== ws && client.readyState === WebSocket.OPEN) {
        client.send(data); // raw forward
      }
    }
  });

  // Clean up on disconnect
  ws.on('close', () => {
    console.log(`❌ Client disconnected: ${clientId}`);
    clients.delete(clientId);

    // Optional: Notify others that a user left
    broadcastExcept(clientId, {
      type: "user-left",
      userId: clientId
    });
  });
});

// Helper function: broadcast to all except one
function broadcastExcept(exceptId, msgObj) {
  const msg = JSON.stringify(msgObj);
  for (const [id, client] of clients) {
    if (id !== exceptId && client.readyState === WebSocket.OPEN) {
      client.send(msg);
    }
  }
}
