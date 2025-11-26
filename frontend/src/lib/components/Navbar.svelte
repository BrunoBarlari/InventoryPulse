<script>
  import { authStore, user, isAdmin } from '../stores/auth.js';
  import { notifications } from '../stores/notifications.js';
  import { theme } from '../stores/theme.js';
  import ConnectionStatus from './ConnectionStatus.svelte';

  function handleLogout() {
    authStore.logout();
    notifications.success('Logged out successfully');
  }
</script>

<nav class="navbar">
  <div class="navbar-inner">
    <div class="navbar-brand">
      <div class="logo">
        <svg width="24" height="24" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
          <rect x="4" y="8" width="32" height="24" rx="2" stroke="currentColor" stroke-width="2.5"/>
          <path d="M4 14H36" stroke="currentColor" stroke-width="2.5"/>
          <path d="M14 14V32" stroke="currentColor" stroke-width="2.5"/>
          <circle cx="9" cy="23" r="2" fill="currentColor"/>
          <circle cx="9" cy="28" r="2" fill="currentColor"/>
        </svg>
      </div>
      <span class="brand-text">InventoryPulse</span>
    </div>

    <div class="navbar-end">
      <ConnectionStatus />

      <!-- Theme Toggle -->
      <button
        class="theme-toggle-btn"
        on:click={() => theme.toggle()}
        aria-label="Toggle dark mode"
        title={$theme === 'light' ? 'Switch to dark mode' : 'Switch to light mode'}
      >
        {#if $theme === 'light'}
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
          </svg>
        {:else}
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="5"/>
            <line x1="12" y1="1" x2="12" y2="3"/>
            <line x1="12" y1="21" x2="12" y2="23"/>
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
            <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
            <line x1="1" y1="12" x2="3" y2="12"/>
            <line x1="21" y1="12" x2="23" y2="12"/>
            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
            <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
          </svg>
        {/if}
      </button>

      {#if $user}
        <div class="divider-vertical"></div>

        <div class="user-info">
          <div class="user-avatar">
            {$user.email.charAt(0).toUpperCase()}
          </div>
          <div class="user-details">
            <span class="user-email">{$user.email}</span>
            <span class="badge {$isAdmin ? 'badge-info' : 'badge-neutral'}">
              {$user.role}
            </span>
          </div>
        </div>

        <button class="btn btn-ghost btn-sm" on:click={handleLogout}>
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
            <polyline points="16 17 21 12 16 7"/>
            <line x1="21" y1="12" x2="9" y2="12"/>
          </svg>
          <span class="logout-text">Log out</span>
        </button>
      {/if}
    </div>
  </div>
</nav>

<style>
  .navbar {
    background: var(--bg-secondary);
    border-bottom: 1px solid var(--border-color);
    position: sticky;
    top: 0;
    z-index: 50;
  }

  .navbar-inner {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 24px;
    max-width: 1400px;
    margin: 0 auto;
  }

  .navbar-brand {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .logo {
    width: 40px;
    height: 40px;
    background: var(--accent-light);
    border-radius: var(--radius-md);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--accent-primary);
  }

  .brand-text {
    font-family: var(--font-display);
    font-size: 1.25rem;
    color: var(--text-primary);
  }

  .navbar-end {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .divider-vertical {
    width: 1px;
    height: 24px;
    background: var(--border-color);
  }

  .theme-toggle-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    background: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
    color: var(--text-secondary);
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .theme-toggle-btn:hover {
    background: var(--bg-primary);
    color: var(--accent-primary);
    border-color: var(--accent-primary);
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .user-avatar {
    width: 36px;
    height: 36px;
    background: var(--bg-tertiary);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    font-size: 14px;
    color: var(--text-secondary);
  }

  .user-details {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .user-email {
    color: var(--text-primary);
    font-size: 14px;
    font-weight: 500;
  }

  .user-details .badge {
    padding: 2px 8px;
    font-size: 11px;
    width: fit-content;
  }

  .logout-text {
    display: inline;
  }

  @media (max-width: 640px) {
    .navbar-inner {
      padding: 12px 16px;
    }

    .user-details {
      display: none;
    }

    .logout-text {
      display: none;
    }

    .brand-text {
      font-size: 1.1rem;
    }
  }
</style>
