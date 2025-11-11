/**
 * Format date and time in the user's local timezone
 */
export function formatDateTime(isoString?: string): { date: string; time: string } {
  if (!isoString) return { date: '—', time: '—' };
  
  const d = new Date(isoString);
  
  // Check if date is valid
  if (isNaN(d.getTime())) return { date: '—', time: '—' };
  
  // Format date as YYYY-MM-DD in local timezone
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  const date = `${year}-${month}-${day}`;
  
  // Format time as HH:MM in local timezone
  const hours = String(d.getHours()).padStart(2, '0');
  const minutes = String(d.getMinutes()).padStart(2, '0');
  const time = `${hours}:${minutes}`;
  
  return { date, time };
}

/**
 * Format date as YYYY-MM-DD in local timezone
 */
export function formatDate(isoString?: string): string {
  if (!isoString) return '—';
  
  const d = new Date(isoString);
  if (isNaN(d.getTime())) return '—';
  
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  
  return `${year}-${month}-${day}`;
}

/**
 * Format time as HH:MM in local timezone
 */
export function formatTime(isoString?: string): string {
  if (!isoString) return '—';
  
  const d = new Date(isoString);
  if (isNaN(d.getTime())) return '—';
  
  const hours = String(d.getHours()).padStart(2, '0');
  const minutes = String(d.getMinutes()).padStart(2, '0');
  
  return `${hours}:${minutes}`;
}

/**
 * Format date and time range in local timezone
 */
export function formatDateTimeRange(startIso?: string, endIso?: string): string {
  if (!startIso || !endIso) return '—';
  
  const start = new Date(startIso);
  const end = new Date(endIso);
  
  if (isNaN(start.getTime()) || isNaN(end.getTime())) return '—';
  
  const date = formatDate(startIso);
  const startTime = formatTime(startIso);
  const endTime = formatTime(endIso);
  
  return `${date} ${startTime} - ${endTime}`;
}

/**
 * Get locale-aware date string (for display in readable format)
 */
export function formatLocaleDate(isoString?: string, options?: Intl.DateTimeFormatOptions): string {
  if (!isoString) return '—';
  
  const d = new Date(isoString);
  if (isNaN(d.getTime())) return '—';
  
  return d.toLocaleDateString(undefined, options || {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
}

/**
 * Get locale-aware time string (for display in readable format)
 */
export function formatLocaleTime(isoString?: string, options?: Intl.DateTimeFormatOptions): string {
  if (!isoString) return '—';
  
  const d = new Date(isoString);
  if (isNaN(d.getTime())) return '—';
  
  return d.toLocaleTimeString(undefined, options || {
    hour: '2-digit',
    minute: '2-digit'
  });
}

/**
 * Get locale-aware date and time string
 */
export function formatLocaleDateTime(isoString?: string): string {
  if (!isoString) return '—';
  
  const d = new Date(isoString);
  if (isNaN(d.getTime())) return '—';
  
  return d.toLocaleString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
}
