const WebSocket = require('ws');
const { v4: uuidv4 } = require('uuid');

const wss = new WebSocket.Server({ port: 8001});
console.log('✅ WebSocket server running at ws://localhost:8001');

const clients = new Map(); // clientId => ws

wss.on('connection', (ws) => {
  const clientId = uuidv4();
  clients.set(clientId, ws);

  console.log(`🔗 Client connected: ${clientId}`);

  // Notify the new client of its ID
  ws.send(JSON.stringify({ type: "welcome", userId: clientId }));

  // Notify others of the new peer
  for (const [id, otherWs] of clients) {
    if (id !== clientId && otherWs.readyState === WebSocket.OPEN) {
      otherWs.send(JSON.stringify({ type: "user-joined", userId: clientId }));
      ws.send(JSON.stringify({ type: "user-joined", userId: id })); // tell new client about existing users
    }
  }

  // Handle messages
  ws.on('message', (data) => {
    let msg;
    try {
      msg = JSON.parse(data);
    } catch (err) {
      console.error('❌ Invalid JSON from client:', err);
      return;
    }

    const { type, target, from } = msg;

    if (type === 'offer' || type === 'answer' || type === 'ice-candidate') {
      const targetWs = clients.get(target);
      if (targetWs && targetWs.readyState === WebSocket.OPEN) {
        targetWs.send(JSON.stringify(msg));
      }
    } else {
      console.log(`📨 Unknown message type: ${type}`);
    }
  });

  ws.on('close', () => {
    console.log(`❌ Client disconnected: ${clientId}`);
    clients.delete(clientId);
  });
});
