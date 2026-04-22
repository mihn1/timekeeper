import { writable } from 'svelte/store';
import { GetRerunJobStatus } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/runtime/runtime';

const defaultStatus = {
  state: 'idle',
  startDate: '',
  endDate: '',
  totalEvents: 0,
  processedEvents: 0,
  errorMessage: '',
  startedAt: '',
  completedAt: '',
  maxRangeDays: 7,
};

export const rerunStatus = writable(defaultStatus);

GetRerunJobStatus()
  .then(s => { if (s) rerunStatus.set({ ...defaultStatus, ...s }); })
  .catch(() => {});

EventsOn('timekeeper:rerun-status', (s) => {
  if (s) rerunStatus.set({ ...defaultStatus, ...s });
});
