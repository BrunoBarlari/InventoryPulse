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
    border-radius: 100px;
    background: #FEF2F2;
    font-size: 12px;
    font-weight: 500;
    border: 1px solid #FECACA;
  }

  .connection-status.connected {
    background: #ECFDF5;
    border-color: #A7F3D0;
  }

  .status-dot {
    width: 6px;
    height: 6px;
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
    50% { opacity: 0.4; }
  }
</style>
