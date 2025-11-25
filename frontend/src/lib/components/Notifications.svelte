<script>
  import { notifications } from '../stores/notifications.js';

  const icons = {
    success: '✓',
    error: '✕',
    warning: '⚠',
    info: 'ℹ',
  };
</script>

<div class="notifications-container">
  {#each $notifications as notification (notification.id)}
    <div
      class="notification glass {notification.type}"
      role="alert"
    >
      <span class="notification-icon">{icons[notification.type]}</span>
      <span class="notification-message">{notification.message}</span>
      <button
        class="notification-close"
        on:click={() => notifications.remove(notification.id)}
      >
        ✕
      </button>
    </div>
  {/each}
</div>

<style>
  .notifications-container {
    position: fixed;
    top: 100px;
    right: 24px;
    display: flex;
    flex-direction: column;
    gap: 12px;
    z-index: 200;
    max-width: 400px;
  }

  .notification {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px 20px;
    animation: slideInRight 0.3s ease-out;
  }

  @keyframes slideInRight {
    from {
      opacity: 0;
      transform: translateX(100px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  .notification-icon {
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    font-size: 12px;
    font-weight: bold;
  }

  .notification.success .notification-icon {
    background: rgba(16, 185, 129, 0.2);
    color: var(--success);
  }

  .notification.error .notification-icon {
    background: rgba(239, 68, 68, 0.2);
    color: var(--error);
  }

  .notification.warning .notification-icon {
    background: rgba(245, 158, 11, 0.2);
    color: var(--warning);
  }

  .notification.info .notification-icon {
    background: rgba(59, 130, 246, 0.2);
    color: var(--info);
  }

  .notification-message {
    flex: 1;
    font-size: 14px;
  }

  .notification-close {
    background: none;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    padding: 4px;
    font-size: 14px;
    transition: color 0.2s;
  }

  .notification-close:hover {
    color: var(--text-primary);
  }
</style>

