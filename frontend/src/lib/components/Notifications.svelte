<script>
  import { notifications } from '../stores/notifications.js';
</script>

<div class="notifications-container">
  {#each $notifications as notification (notification.id)}
    <div
      class="notification {notification.type}"
      role="alert"
    >
      <span class="notification-icon">
        {#if notification.type === 'success'}
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <polyline points="20 6 9 17 4 12"/>
          </svg>
        {:else if notification.type === 'error'}
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        {:else if notification.type === 'warning'}
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="M12 9v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        {:else}
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="16" x2="12" y2="12"/>
            <line x1="12" y1="8" x2="12.01" y2="8"/>
          </svg>
        {/if}
      </span>
      <span class="notification-message">{notification.message}</span>
      <button
        class="notification-close"
        on:click={() => notifications.remove(notification.id)}
        aria-label="Dismiss notification"
      >
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"/>
          <line x1="6" y1="6" x2="18" y2="18"/>
        </svg>
      </button>
    </div>
  {/each}
</div>

<style>
  .notifications-container {
    position: fixed;
    top: 80px;
    right: 24px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    z-index: 200;
    max-width: 380px;
  }

  .notification {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 14px 16px;
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-lg);
    animation: slideInRight 0.25s ease-out;
  }

  @keyframes slideInRight {
    from {
      opacity: 0;
      transform: translateX(50px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  .notification-icon {
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    flex-shrink: 0;
  }

  .notification.success {
    border-left: 3px solid var(--success);
  }

  .notification.success .notification-icon {
    background: #ECFDF5;
    color: var(--success);
  }

  .notification.error {
    border-left: 3px solid var(--error);
  }

  .notification.error .notification-icon {
    background: #FEF2F2;
    color: var(--error);
  }

  .notification.warning {
    border-left: 3px solid var(--warning);
  }

  .notification.warning .notification-icon {
    background: #FFFBEB;
    color: var(--warning);
  }

  .notification.info {
    border-left: 3px solid var(--info);
  }

  .notification.info .notification-icon {
    background: #F0F9FF;
    color: var(--info);
  }

  .notification-message {
    flex: 1;
    font-size: 14px;
    color: var(--text-primary);
    line-height: 1.4;
  }

  .notification-close {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    background: transparent;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    border-radius: var(--radius-sm);
    transition: all 0.15s ease;
    flex-shrink: 0;
  }

  .notification-close:hover {
    background: var(--bg-tertiary);
    color: var(--text-primary);
  }

  @media (max-width: 480px) {
    .notifications-container {
      left: 16px;
      right: 16px;
      max-width: none;
    }
  }
</style>
