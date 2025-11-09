// rule.js
import { WF_STATUSES, SE_CODE_USER_TYPE, SE_ACCNT } from "./details.js";
import { escapeHtml, nullable } from "../helpers.js";

function opt(v, t) { return `<option value="${v}">${escapeHtml(t)}</option>`; }
function buildOptions(arr, idKey, textKey, extraText = null) {
  return arr.map(o => {
    const id = o[idKey];
    const text = extraText ? extraText(o) : o[textKey];
    return opt(id, text);
  }).join('');
}

window.createRule = createRule;
export function createRule(e, workflowID, version) {
  if (e && typeof e.preventDefault === 'function') e.preventDefault();

  const form = document.getElementById('add-rule-form');
  const formData = new FormData(form);
  const data = Object.fromEntries(formData.entries());
  data.workflow_id = workflowID;
  data.version = version;

  fetch('/add_rule', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  .then(r => { if (!r.ok) throw new Error('Failed to create rule'); })
  .then(() => {
    const wfid = new URLSearchParams(window.location.search).get('workflow');
    if (wfid) window.viewDetails(wfid);
  })
  .catch(err => console.error('Error creating rule:', err));
}

window.editRule = editRule;
export function editRule(rule_id, from_status_id, from_status_name, to_status_id, to_status_name, se_code_user_type_id, user_type_en, se_accnt_id, accnt_en, action_button, action_function, workflowID, version, is_condition) {
  const row = document.getElementById(`rule_row_${rule_id}`);
  if (!row) return;

  const formId = `edit_rule_form_${rule_id}`;
  const fromSelId = `from_status_sel_${rule_id}`;
  const toSelId   = `to_status_sel_${rule_id}`;
  const utSelId   = `user_type_sel_${rule_id}`;
  const accSelId  = `account_sel_${rule_id}`;

  row.innerHTML = `
    <form id="${formId}" style="display:none"></form>

    <td>${rule_id}</td>

    <td>
      <select form="${formId}" id="${fromSelId}" name="from_status_id" required>
        <option value="">--</option>
        ${buildOptions(WF_STATUSES,'status_id','status_name',o=>`${o.status_name} (ID:${o.status_id})`)}
      </select>
    </td>

    <td>
      <select form="${formId}" id="${toSelId}" name="to_status_id" required>
        <option value="">--</option>
        ${buildOptions(WF_STATUSES,'status_id','status_name',o=>`${o.status_name} (ID:${o.status_id})`)}
      </select>
    </td>

    <td>
      <select form="${formId}" id="${utSelId}" name="se_code_user_type_id">
        <option value="">--</option>
        ${buildOptions(SE_CODE_USER_TYPE,'se_code_user_type_id','descr_en',o=>`${o.descr_en} (ID:${o.se_code_user_type_id})`)}
      </select>
    </td>

    <td>
      <select form="${formId}" id="${accSelId}" name="se_accnt_id">
        <option value="">--</option>
        ${buildOptions(SE_ACCNT,'se_accnt_id','descr_en',o=>`${o.descr_en} (ID:${o.se_accnt_id})`)}
      </select>
    </td>

    <td>
      <input form="${formId}" type="text" name="action_button" value="${escapeHtml(action_button||'')}" required>
    </td>

    <td>
      <input form="${formId}" type="text" name="action_function" value="${escapeHtml(action_function||'')}">
    </td>

    <td>
      <button class="icon-btn save" type="button"
              onclick="event.stopPropagation(); saveRule(${rule_id}, ${workflowID}, ${version})"
              title="Save">
        <i class="fa-solid fa-check"></i>
      </button>
      <button class="icon-btn cancel" type="button"
              onclick="event.stopPropagation(); cancelRuleEdit(${rule_id}, '${from_status_id}', '${escapeHtml(from_status_name||'')}', '${to_status_id}', '${escapeHtml(to_status_name||'')}', '${se_code_user_type_id}', '${escapeHtml(user_type_en||'')}', '${se_accnt_id}', '${escapeHtml(accnt_en||'')}', '${escapeHtml(action_button||'')}', '${escapeHtml(action_function||'')}', ${workflowID}, ${version}, ${is_condition})"
              title="Cancel">
        <i class="fa-solid fa-xmark"></i>
      </button>
    </td>
  `;

  const tsFrom = new TomSelect('#'+fromSelId, { create:false, closeAfterSelect:true, allowEmptyOption:true });
  const tsTo   = new TomSelect('#'+toSelId,   { create:false, closeAfterSelect:true, allowEmptyOption:true });
  const tsUT   = new TomSelect('#'+utSelId,   { create:false, closeAfterSelect:true, allowEmptyOption:true });
  const tsAcc  = new TomSelect('#'+accSelId,  { create:false, closeAfterSelect:true, allowEmptyOption:true });

  if (from_status_id) tsFrom.setValue(String(from_status_id), true);
  if (to_status_id)   tsTo.setValue(String(to_status_id), true);
  if (se_code_user_type_id) tsUT.setValue(String(se_code_user_type_id), true);
  if (se_accnt_id)          tsAcc.setValue(String(se_accnt_id), true);

  // account → user type auto-fill
  tsAcc.on('change', (val)=>{
    const acct = SE_ACCNT.find(a => String(a.se_accnt_id) === String(val));
    if (!acct) return;
    const ut = SE_CODE_USER_TYPE.find(u => String(u.se_code_user_type_id) === String(acct.se_code_user_type_id));
    if (ut) tsUT.setValue(String(ut.se_code_user_type_id), true);
  });
}

window.deleteRule = deleteRule;
export function deleteRule(rule_id, workflowID, version) {
  if (!confirm("Are you sure you want to delete this rule?")) return;
  fetch(`/delete_rule?rule_id=${rule_id}&workflow_id=${workflowID}&version=${version}`)
    .then(r => { if (!r.ok) throw new Error('Failed'); })
    .then(() => {
      const wfid = new URLSearchParams(window.location.search).get('workflow');
      if (wfid) window.viewDetails(wfid);
    })
    .catch(err => console.error("Error deleting rule:", err));
}

window.saveRule = saveRule;
export function saveRule(rule_id, workflowID, version){
  const form = document.getElementById('edit_rule_form_' + rule_id);
  if (!form) return console.error('edit_rule_form not found');

  const formData = new FormData(form);
  const params = new URLSearchParams({
    workflow_id: workflowID,
    version: version,
    rule_id: rule_id,
    from_status_id: formData.get('from_status_id') || '',
    to_status_id: formData.get('to_status_id') || '',
    se_code_user_type_id: formData.get('se_code_user_type_id') || '',
    se_accnt_id: formData.get('se_accnt_id') || '',
    action_button: formData.get('action_button') || '',
    action_function: formData.get('action_function') || '',
  });

  fetch(`/edit_rule?${params.toString()}`, { method:'GET', headers:{'Cache-Control':'no-cache'} })
    .then(r => { if (!r.ok) throw new Error("Failed to update rule"); return r.json(); })
    .then(() => {
      const wfid = new URLSearchParams(window.location.search).get('workflow');
      if (wfid) window.viewDetails(wfid);
    })
    .catch(err => console.error("Error updating rule:", err));
}

window.cancelRuleEdit = cancelRuleEdit;
export function cancelRuleEdit(rule_id, from_status_id, from_status_name, to_status_id, to_status_name, se_code_user_type_id, user_type_en, se_accnt_id, accnt_en, action_button, action_function, workflowID, version, is_condition){
  const row = document.getElementById(`rule_row_${rule_id}`);
  if (!row) return;

  row.innerHTML = `
    <td>${rule_id}</td>
    <td>${nullable(from_status_name)}</td>
    <td>${nullable(to_status_name)}</td>
    <td>${nullable(user_type_en)}</td>
    <td>${nullable(accnt_en)}</td>
    <td>${escapeHtml(action_button || '')}</td>
    <td>${escapeHtml(action_function || '')}</td>
    <td>
      <button class="icon-btn edit" type="button"
        onclick="event.stopPropagation(); editRule(${rule_id}, ${from_status_id}, '${escapeHtml(from_status_name||'')}', ${to_status_id}, '${escapeHtml(to_status_name||'')}', ${se_code_user_type_id}, '${escapeHtml(user_type_en||'')}', ${se_accnt_id}, '${escapeHtml(accnt_en||'')}', '${escapeHtml(action_button||'')}', '${escapeHtml(action_function||'')}', ${workflowID}, ${version})"
        title="Edit">
        <i class="fa-solid fa-pen"></i>
      </button>

      <button class="icon-btn delete" type="button"
        onclick="event.stopPropagation(); deleteRule(${rule_id}, ${workflowID}, ${version})"
        title="Delete">
        <i class="fa-solid fa-trash"></i>
      </button>

      <button class="icon-btn" onclick="event.stopPropagation(); toggleRuleConditions(${rule_id})" title="Toggle conditions">
        <i class="fa-solid fa-plus" id="icon_plus_${rule_id}" style="transform: rotate(${is_condition ? '45deg' : '0deg'}); transition: transform 0.2s;"></i>
      </button>
    </td>
  `;
}
