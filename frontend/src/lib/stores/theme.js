import { writable } from 'svelte/store';

// Check for saved theme or system preference
function getInitialTheme() {
  if (typeof window === 'undefined') return 'light';

  const saved = localStorage.getItem('theme');
  if (saved) return saved;

  // Check system preference
  if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    return 'dark';
  }
  return 'light';
}

function createThemeStore() {
  const { subscribe, set, update } = writable(getInitialTheme());

  return {
    subscribe,
    toggle: () => {
      update(current => {
        const newTheme = current === 'light' ? 'dark' : 'light';
        localStorage.setItem('theme', newTheme);
        document.documentElement.setAttribute('data-theme', newTheme);
        return newTheme;
      });
    },
    init: () => {
      const theme = getInitialTheme();
      document.documentElement.setAttribute('data-theme', theme);
      set(theme);
    }
  };
}

export const theme = createThemeStore();

