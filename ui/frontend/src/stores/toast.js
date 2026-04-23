import { writable } from 'svelte/store';

export const toasts = writable([]);

let nextId = 1;

export function showToast(message, { type = 'success', duration = 3000 } = {}) {
  const id = nextId++;
  toasts.update(list => [...list, { id, message, type }]);
  if (duration > 0) {
    setTimeout(() => dismissToast(id), duration);
  }
  return id;
}

export function dismissToast(id) {
  toasts.update(list => list.filter(t => t.id !== id));
}
