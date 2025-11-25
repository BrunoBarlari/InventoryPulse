import { writable, get } from 'svelte/store';
import { notifications } from './notifications.js';

const WS_URL = 'ws://localhost:8080/ws';

function createWebSocketStore() {
  const { subscribe, set, update } = writable({
    connected: false,
    lastMessage: null,
  });

  let ws = null;
  let reconnectAttempts = 0;
  const maxReconnectAttempts = 5;
  const reconnectDelay = 3000;

  // Event handlers that will be set by consumers
  let eventHandlers = {};

  function connect() {
    if (ws && ws.readyState === WebSocket.OPEN) {
      return;
    }

    try {
      ws = new WebSocket(WS_URL);

      ws.onopen = () => {
        console.log('WebSocket connected');
        reconnectAttempts = 0;
        update(state => ({ ...state, connected: true }));
        notifications.info('Real-time updates enabled');
      };

      ws.onclose = () => {
        console.log('WebSocket disconnected');
        update(state => ({ ...state, connected: false }));
        
        // Attempt to reconnect
        if (reconnectAttempts < maxReconnectAttempts) {
          reconnectAttempts++;
          console.log(`Reconnecting... attempt ${reconnectAttempts}`);
          setTimeout(connect, reconnectDelay);
        }
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };

      ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          update(state => ({ ...state, lastMessage: message }));
          
          // Call registered event handlers
          if (eventHandlers[message.type]) {
            eventHandlers[message.type].forEach(handler => handler(message.payload));
          }

          // Show notification for events
          showEventNotification(message);
        } catch (err) {
          console.error('Error parsing WebSocket message:', err);
        }
      };
    } catch (err) {
      console.error('Error connecting to WebSocket:', err);
    }
  }

  function disconnect() {
    if (ws) {
      ws.close();
      ws = null;
    }
    update(state => ({ ...state, connected: false }));
  }

  function showEventNotification(message) {
    switch (message.type) {
      case 'product.created':
        notifications.success(`New product: ${message.payload.name}`);
        break;
      case 'product.updated':
        notifications.info(`Product updated: ${message.payload.name}`);
        break;
      case 'product.deleted':
        notifications.warning('Product deleted');
        break;
      case 'stock.updated':
        notifications.info(`Stock updated: ${message.payload.name} → ${message.payload.quantity}`);
        break;
    }
  }

  function on(eventType, handler) {
    if (!eventHandlers[eventType]) {
      eventHandlers[eventType] = [];
    }
    eventHandlers[eventType].push(handler);

    // Return unsubscribe function
    return () => {
      eventHandlers[eventType] = eventHandlers[eventType].filter(h => h !== handler);
    };
  }

  function off(eventType, handler) {
    if (eventHandlers[eventType]) {
      eventHandlers[eventType] = eventHandlers[eventType].filter(h => h !== handler);
    }
  }

  return {
    subscribe,
    connect,
    disconnect,
    on,
    off,
  };
}

export const websocketStore = createWebSocketStore();

