<script>
  import { authStore, user, isAdmin } from '../stores/auth.js';
  import { notifications } from '../stores/notifications.js';
  import ConnectionStatus from './ConnectionStatus.svelte';

  function handleLogout() {
    authStore.logout();
    notifications.success('Logged out successfully');
  }
</script>

<nav class="navbar glass">
  <div class="navbar-brand">
    <span class="logo">📦</span>
    <span class="brand-text">InventoryPulse</span>
  </div>

  <div class="navbar-user">
    <ConnectionStatus />
    {#if $user}
      <div class="user-info">
        <span class="user-email">{$user.email}</span>
        <span class="badge {$isAdmin ? 'badge-info' : 'badge-success'}">
          {$user.role}
        </span>
      </div>
      <button class="btn btn-secondary btn-sm" on:click={handleLogout}>
        Logout
      </button>
    {/if}
  </div>
</nav>

<style>
  .navbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 32px;
    margin: 16px;
    position: sticky;
    top: 16px;
    z-index: 50;
  }

  .navbar-brand {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .logo {
    font-size: 28px;
  }

  .brand-text {
    font-family: var(--font-display);
    font-size: 1.5rem;
    font-weight: 700;
    background: var(--accent-gradient);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .navbar-user {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .user-email {
    color: var(--text-secondary);
    font-size: 14px;
  }
</style>
