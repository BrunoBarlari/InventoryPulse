<script>
  import { onMount } from 'svelte';
  import { authStore, isAuthenticated, isLoading } from './lib/stores/auth.js';
  import Notifications from './lib/components/Notifications.svelte';
  import Login from './routes/Login.svelte';
  import Dashboard from './routes/Dashboard.svelte';

  onMount(() => {
    authStore.init();
  });
</script>

<Notifications />

{#if $isLoading}
  <div class="app-loading">
    <div class="loading-content">
      <div class="loading-logo">
        <svg width="48" height="48" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
          <rect x="4" y="8" width="32" height="24" rx="2" stroke="currentColor" stroke-width="2"/>
          <path d="M4 14H36" stroke="currentColor" stroke-width="2"/>
          <path d="M14 14V32" stroke="currentColor" stroke-width="2"/>
          <circle cx="9" cy="23" r="2" fill="currentColor"/>
          <circle cx="9" cy="28" r="2" fill="currentColor"/>
        </svg>
      </div>
      <div class="spinner spinner-lg"></div>
      <p>Loading InventoryPulse...</p>
    </div>
  </div>
{:else if $isAuthenticated}
  <Dashboard />
{:else}
  <Login on:login={() => {}} />
{/if}

<style>
  .app-loading {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--bg-primary);
  }

  .loading-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
    text-align: center;
  }

  .loading-logo {
    width: 72px;
    height: 72px;
    background: var(--accent-light);
    border-radius: var(--radius-xl);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--accent-primary);
  }

  .loading-content p {
    color: var(--text-muted);
    font-size: 14px;
  }
</style>
