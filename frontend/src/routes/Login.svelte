<script>
  import { authStore } from '../lib/stores/auth.js';
  import { notifications } from '../lib/stores/notifications.js';
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  let email = '';
  let password = '';
  let isLoading = false;
  let error = '';

  async function handleSubmit() {
    if (!email || !password) {
      error = 'Please fill in all fields';
      return;
    }

    isLoading = true;
    error = '';

    try {
      await authStore.login(email, password);
      notifications.success('Welcome back!');
      dispatch('login');
    } catch (err) {
      error = err.message || 'Invalid email or password';
      notifications.error(error);
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="login-page">
  <div class="login-container">
    <div class="login-content animate-fadeIn">
      <!-- Header -->
      <header class="login-header">
        <div class="logo">
          <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect x="4" y="8" width="32" height="24" rx="2" stroke="currentColor" stroke-width="2"/>
            <path d="M4 14H36" stroke="currentColor" stroke-width="2"/>
            <path d="M14 14V32" stroke="currentColor" stroke-width="2"/>
            <circle cx="9" cy="23" r="2" fill="currentColor"/>
            <circle cx="9" cy="28" r="2" fill="currentColor"/>
          </svg>
        </div>
        <h1>InventoryPulse</h1>
        <p class="subtitle">Inventory management, simplified.</p>
      </header>

      <!-- Login Form -->
      <form on:submit|preventDefault={handleSubmit} class="login-form">
        {#if error}
          <div class="error-alert">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
            {error}
          </div>
        {/if}

        <div class="input-group">
          <label for="email">Email address</label>
          <input
            type="email"
            id="email"
            class="input"
            placeholder="you@example.com"
            bind:value={email}
            disabled={isLoading}
            autocomplete="email"
          />
        </div>

        <div class="input-group">
          <label for="password">Password</label>
          <input
            type="password"
            id="password"
            class="input"
            placeholder="Enter your password"
            bind:value={password}
            disabled={isLoading}
            autocomplete="current-password"
          />
        </div>

        <button type="submit" class="btn btn-primary w-full" disabled={isLoading}>
          {#if isLoading}
            <span class="spinner"></span>
            Signing in...
          {:else}
            Sign in
          {/if}
        </button>
      </form>

      <!-- Demo credentials -->
      <div class="demo-info">
        <p class="demo-label">Demo credentials</p>
        <div class="demo-credentials">
          <code>admin@inventorypulse.com</code>
          <span class="separator">/</span>
          <code>admin123</code>
        </div>
      </div>
    </div>

    <!-- Side illustration -->
    <div class="login-illustration">
      <div class="illustration-content">
        <div class="floating-cards">
          <div class="float-card card-1">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
            </svg>
            <span>Products</span>
          </div>
          <div class="float-card card-2">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
            </svg>
            <span>Analytics</span>
          </div>
          <div class="float-card card-3">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M13 10V3L4 14h7v7l9-11h-7z"/>
            </svg>
            <span>Real-time</span>
          </div>
        </div>
        <div class="grid-pattern"></div>
      </div>
    </div>
  </div>
</div>

<style>
  .login-page {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 24px;
    background: var(--bg-primary);
  }

  .login-container {
    display: grid;
    grid-template-columns: 1fr 1fr;
    max-width: 1000px;
    width: 100%;
    background: var(--bg-secondary);
    border-radius: var(--radius-xl);
    border: 1px solid var(--border-color);
    box-shadow: var(--shadow-lg);
    overflow: hidden;
  }

  @media (max-width: 768px) {
    .login-container {
      grid-template-columns: 1fr;
    }
    .login-illustration {
      display: none;
    }
  }

  .login-content {
    padding: 48px 40px;
    display: flex;
    flex-direction: column;
  }

  .login-header {
    margin-bottom: 40px;
  }

  .logo {
    width: 56px;
    height: 56px;
    background: var(--accent-light);
    border-radius: var(--radius-lg);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--accent-primary);
    margin-bottom: 24px;
  }

  .login-header h1 {
    font-size: 1.75rem;
    margin-bottom: 8px;
    color: var(--text-primary);
  }

  .subtitle {
    color: var(--text-muted);
    font-size: 15px;
  }

  .login-form {
    display: flex;
    flex-direction: column;
    gap: 20px;
    flex: 1;
  }

  .error-alert {
    display: flex;
    align-items: center;
    gap: 10px;
    background: #FEF2F2;
    border: 1px solid #FECACA;
    color: var(--error);
    padding: 12px 14px;
    border-radius: var(--radius-md);
    font-size: 14px;
  }

  .error-alert svg {
    flex-shrink: 0;
  }

  .demo-info {
    margin-top: auto;
    padding-top: 32px;
    border-top: 1px solid var(--border-color);
  }

  .demo-label {
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--text-muted);
    margin-bottom: 8px;
  }

  .demo-credentials {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }

  .demo-credentials code {
    background: var(--bg-tertiary);
    padding: 4px 10px;
    border-radius: 4px;
    font-size: 13px;
    color: var(--text-secondary);
  }

  .separator {
    color: var(--text-muted);
  }

  /* Illustration side */
  .login-illustration {
    background: linear-gradient(135deg, var(--accent-light) 0%, #F0FDF4 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    overflow: hidden;
  }

  .illustration-content {
    position: relative;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .floating-cards {
    position: relative;
    z-index: 2;
  }

  .float-card {
    display: flex;
    align-items: center;
    gap: 12px;
    background: var(--bg-secondary);
    padding: 16px 20px;
    border-radius: var(--radius-lg);
    border: 1px solid var(--border-color);
    box-shadow: var(--shadow-md);
    position: absolute;
    white-space: nowrap;
    font-weight: 500;
    font-size: 14px;
    color: var(--text-primary);
    animation: float 6s ease-in-out infinite;
  }

  .float-card svg {
    color: var(--accent-primary);
  }

  .card-1 {
    top: 60px;
    left: 30px;
    animation-delay: 0s;
  }

  .card-2 {
    top: 140px;
    right: 40px;
    animation-delay: -2s;
  }

  .card-3 {
    bottom: 80px;
    left: 60px;
    animation-delay: -4s;
  }

  @keyframes float {
    0%, 100% {
      transform: translateY(0);
    }
    50% {
      transform: translateY(-10px);
    }
  }

  .grid-pattern {
    position: absolute;
    inset: 0;
    background-image:
      linear-gradient(var(--border-color) 1px, transparent 1px),
      linear-gradient(90deg, var(--border-color) 1px, transparent 1px);
    background-size: 40px 40px;
    opacity: 0.5;
  }
</style>
