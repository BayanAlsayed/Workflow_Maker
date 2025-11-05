
// helpers
export function nullable(v) {
  return v === null || v === undefined ? '' : escapeHtml(String(v));
}

export function escapeHtml(s) {
  return s.replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));
}

// select option builders
export function opt(v, t) { return `<option value="${v}">${escapeHtml(t)}</option>`; }
export function buildOptions(arr, idKey, textKey, extraText = null) {
  return arr.map(o => {
    const id = o[idKey];
    const text = extraText ? extraText(o) : o[textKey];
    return opt(id, text);
  }).join('');
}

export function initTomSelectSingle(selector, onChange) {
  const el = document.querySelector(selector);
  if (!el) return null;
  const ts = new TomSelect(el, { create:false, maxOptions:1000, closeAfterSelect:true, allowEmptyOption:true });
  if (onChange) ts.on('change', onChange);
  return ts;
}