
// helpers
export function nullable(v) {
  return v === null || v === undefined ? '' : escapeHtml(String(v));
}

export function escapeHtml(s) {
  return s.replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));
}