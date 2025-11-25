import { writable, derived } from 'svelte/store';
import { auth as authAPI } from '../api.js';

// Create stores
function createAuthStore() {
  const { subscribe, set, update } = writable({
    user: null,
    isAuthenticated: false,
    isLoading: true,
  });

  return {
    subscribe,

    // Initialize from localStorage
    init() {
      const isAuth = authAPI.isAuthenticated();
      const user = authAPI.getUser();

      set({
        user,
        isAuthenticated: isAuth,
        isLoading: false,
      });
    },

    // Login
    async login(email, password) {
      const { user } = await authAPI.login(email, password);
      set({
        user,
        isAuthenticated: true,
        isLoading: false,
      });
      return user;
    },

    // Logout
    logout() {
      authAPI.logout();
      set({
        user: null,
        isAuthenticated: false,
        isLoading: false,
      });
    },

    // Update user data
    setUser(user) {
      update(state => ({ ...state, user }));
    },
  };
}

export const authStore = createAuthStore();

// Derived stores for convenience
export const isAuthenticated = derived(authStore, $auth => $auth.isAuthenticated);
export const user = derived(authStore, $auth => $auth.user);
export const isAdmin = derived(authStore, $auth => $auth.user?.role === 'admin');
export const isLoading = derived(authStore, $auth => $auth.isLoading);

