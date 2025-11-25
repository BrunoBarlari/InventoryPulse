<script>
  import { websocketStore } from '../stores/websocket.js';
</script>

<div class="connection-status" class:connected={$websocketStore.connected}>
  <span class="status-dot"></span>
  <span class="status-text">
    {$websocketStore.connected ? 'Live' : 'Offline'}
  </span>
</div>

<style>
  .connection-status {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 12px;
    border-radius: 20px;
    background: rgba(239, 68, 68, 0.2);
    font-size: 12px;
    font-weight: 500;
  }

  .connection-status.connected {
    background: rgba(16, 185, 129, 0.2);
  }

  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--error);
    animation: pulse 2s infinite;
  }

  .connected .status-dot {
    background: var(--success);
    animation: none;
  }

  .status-text {
    color: var(--error);
  }

  .connected .status-text {
    color: var(--success);
  }

  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.5; }
  }
</style>

