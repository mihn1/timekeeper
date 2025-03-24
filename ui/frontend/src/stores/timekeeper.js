import { writable } from 'svelte/store';

export const trackingEnabled = writable(false);
export const refreshData = writable(Date.now());