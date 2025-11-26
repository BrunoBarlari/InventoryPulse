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
    gap: 8px;
    padding: 6px 14px;
    border-radius: 100px;
    background: var(--error-bg);
    font-size: 12px;
    font-weight: 500;
    border: 1px solid transparent;
    transition: all 0.3s ease;
  }

  .connection-status.connected {
    background: var(--success-bg);
    border-color: rgba(5, 150, 105, 0.2);
  }

  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--error);
    position: relative;
  }

  /* Pulsing animation for connected state */
  .connected .status-dot {
    background: var(--success);
    animation: pulse-dot 2s ease-in-out infinite;
  }

  .connected .status-dot::before {
    content: '';
    position: absolute;
    inset: -3px;
    border-radius: 50%;
    background: var(--success);
    animation: pulse-ring 2s ease-in-out infinite;
  }

  @keyframes pulse-dot {
    0%, 100% {
      transform: scale(1);
      opacity: 1;
    }
    50% {
      transform: scale(1.1);
      opacity: 0.9;
    }
  }

  @keyframes pulse-ring {
    0% {
      transform: scale(1);
      opacity: 0.4;
    }
    50% {
      transform: scale(1.8);
      opacity: 0;
    }
    100% {
      transform: scale(1);
      opacity: 0;
    }
  }

  /* Disconnected animation */
  .status-dot:not(.connected .status-dot) {
    animation: blink 1.5s ease-in-out infinite;
  }

  @keyframes blink {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.3; }
  }

  .status-text {
    color: var(--error);
    letter-spacing: 0.02em;
  }

  .connected .status-text {
    color: var(--success);
  }
</style>
