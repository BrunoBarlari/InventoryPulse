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
    <div class="spinner"></div>
    <p>Loading...</p>
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
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 16px;
    color: var(--text-muted);
  }
</style>
