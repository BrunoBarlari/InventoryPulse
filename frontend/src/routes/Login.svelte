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
      notifications.success('Welcome back! 🎉');
      dispatch('login');
    } catch (err) {
      error = err.message || 'Invalid email or password';
      notifications.error(error);
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="login-container">
  <div class="login-card glass animate-fadeIn">
    <div class="login-header">
      <span class="login-icon">📦</span>
      <h1>InventoryPulse</h1>
      <p class="text-secondary">Manage your inventory in real-time</p>
    </div>

    <form on:submit|preventDefault={handleSubmit}>
      {#if error}
        <div class="error-message">
          {error}
        </div>
      {/if}

      <div class="input-group">
        <label for="email">Email</label>
        <input
          type="email"
          id="email"
          class="input"
          placeholder="admin@inventorypulse.com"
          bind:value={email}
          disabled={isLoading}
        />
      </div>

      <div class="input-group">
        <label for="password">Password</label>
        <input
          type="password"
          id="password"
          class="input"
          placeholder="••••••••"
          bind:value={password}
          disabled={isLoading}
        />
      </div>

      <button type="submit" class="btn btn-primary w-full" disabled={isLoading}>
        {#if isLoading}
          <span class="spinner"></span>
          Signing in...
        {:else}
          Sign In
        {/if}
      </button>
    </form>

    <div class="login-footer">
      <p class="text-muted">
        Demo credentials: <br />
        <code>admin@inventorypulse.com</code> / <code>admin123</code>
      </p>
    </div>
  </div>

  <div class="login-decoration">
    <div class="blob blob-1"></div>
    <div class="blob blob-2"></div>
    <div class="blob blob-3"></div>
  </div>
</div>

<style>
  .login-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 24px;
    position: relative;
    overflow: hidden;
  }

  .login-card {
    width: 100%;
    max-width: 420px;
    padding: 48px 40px;
    position: relative;
    z-index: 10;
  }

  .login-header {
    text-align: center;
    margin-bottom: 40px;
  }

  .login-icon {
    font-size: 64px;
    display: block;
    margin-bottom: 16px;
  }

  .login-header h1 {
    font-size: 2rem;
    margin-bottom: 8px;
    background: var(--accent-gradient);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .error-message {
    background: rgba(239, 68, 68, 0.15);
    border: 1px solid rgba(239, 68, 68, 0.3);
    color: var(--error);
    padding: 12px 16px;
    border-radius: var(--radius-md);
    font-size: 14px;
    text-align: center;
  }

  .login-footer {
    margin-top: 32px;
    text-align: center;
    font-size: 13px;
  }

  .login-footer code {
    background: rgba(255, 255, 255, 0.1);
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 12px;
  }

  /* Decorative blobs */
  .login-decoration {
    position: absolute;
    inset: 0;
    overflow: hidden;
    pointer-events: none;
  }

  .blob {
    position: absolute;
    border-radius: 50%;
    filter: blur(80px);
    opacity: 0.5;
    animation: float 20s infinite ease-in-out;
  }

  .blob-1 {
    width: 400px;
    height: 400px;
    background: var(--accent-primary);
    top: -100px;
    right: -100px;
    animation-delay: 0s;
  }

  .blob-2 {
    width: 300px;
    height: 300px;
    background: var(--accent-secondary);
    bottom: -50px;
    left: -50px;
    animation-delay: -7s;
  }

  .blob-3 {
    width: 250px;
    height: 250px;
    background: #ec4899;
    top: 50%;
    left: 30%;
    animation-delay: -14s;
  }

  @keyframes float {
    0%, 100% {
      transform: translate(0, 0) scale(1);
    }
    25% {
      transform: translate(20px, -30px) scale(1.05);
    }
    50% {
      transform: translate(-20px, 20px) scale(0.95);
    }
    75% {
      transform: translate(30px, 10px) scale(1.02);
    }
  }
</style>

