import { writable } from 'svelte/store';

// Check if dark mode is preferred
const prefersDarkMode = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
const initialTheme = localStorage.getItem('theme') || (prefersDarkMode ? 'dark' : 'light');

// Apply theme to document immediately
function applyTheme(theme) {
  document.documentElement.setAttribute('data-theme', theme);
  localStorage.setItem('theme', theme);
}

// Initialize the theme store
export const theme = writable(initialTheme);

// Apply initial theme
applyTheme(initialTheme);

// Subscribe to theme changes
theme.subscribe(value => {
  applyTheme(value);
});

// Toggle theme function
export function toggleTheme() {
  theme.update(t => t === 'light' ? 'dark' : 'light');
}
