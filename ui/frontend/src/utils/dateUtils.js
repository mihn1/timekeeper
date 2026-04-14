/**
 * Returns "YYYY-MM-DD" for today in the given IANA timezone.
 * Uses en-CA locale which always produces YYYY-MM-DD output.
 */
export function todayInTz(tz) {
  try {
    return new Intl.DateTimeFormat('en-CA', { timeZone: tz }).format(new Date());
  } catch {
    return localDateStr(new Date());
  }
}

/**
 * Shifts a "YYYY-MM-DD" string by `days` calendar days.
 * Timezone-independent: operates entirely on the calendar date.
 */
export function shiftDateStr(dateStr, days) {
  const d = new Date(dateStr + 'T00:00:00');
  d.setDate(d.getDate() + days);
  return localDateStr(d);
}

/**
 * Formats a "YYYY-MM-DD" string for user display in the given timezone.
 * Uses noon to avoid DST-boundary edge cases with midnight.
 */
export function formatDateDisplay(dateStr, tz) {
  try {
    const d = new Date(dateStr + 'T12:00:00');
    return new Intl.DateTimeFormat(undefined, {
      timeZone: tz,
      weekday: 'short',
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    }).format(d);
  } catch {
    return dateStr;
  }
}

function localDateStr(d) {
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${y}-${m}-${day}`;
}
