// status.js
import { escapeHtml, nullable } from '../helpers.js';
import { ED_CODE_STATUS_CAT, ED_CODE_STATUS, GS_CODE_REQ_STATUS } from './details.js';

function opt(v, t) { return `<option value="${v}">${escapeHtml(t)}</option>`; }
function buildOptions(arr, idKey, textKey, extraText = null) {
  return arr.map(o => {
    const id = o[idKey];
    const text = extraText ? extraText(o) : o[textKey];
    return opt(id, text);
  }).join('');
}

window.createStatus = createStatus;
export function createStatus(e, workflowID, version) {
  if (e && typeof e.preventDefault === 'function') e.preventDefault();

  const form = document.getElementById('add-status-form');
  if (!form) return console.error('add-status-form not found');

  const formData = new FormData(form);
  const edId = formData.get('ed_code_status_id');
  const gsId = formData.get('gs_code_req_status_id');

  if ((!edId && !gsId) || (edId && gsId)) {
    alert("Please select either an ED Code Status or a GS Code Request Status.");
    return;
  }

  const data = Object.fromEntries(formData.entries());
  data.workflow_id = workflowID;
  data.version =  version

  fetch('/add_status', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  .then(r => { if (!r.ok) throw new Error('Failed to create status'); })
  .then(() => {
    const wfid = new URLSearchParams(window.location.search).get('workflow');
    if (wfid) window.viewDetails(wfid);
  })
  .catch(err => console.error('Error creating status:', err));
}

window.editStatus = editStatus;
export function editStatus(statusID, statusName, edID, ED_DESCR_EN, gsID, GS_DESCR_EN, isTerminal, successPath, workflowID, version) {
  const row = document.getElementById(`status_row_${statusID}`);
  if (!row) return;

  const formId     = `edit_status_form_${statusID}`;
  const edCatSelId = `ed_cat_sel_${statusID}`;
  const edCodeSelId= `ed_code_sel_${statusID}`;
  const gsSelId    = `gs_sel_${statusID}`;

  // derive cat from edID if present
  const edCatVal = (function() {
    if (!edID) return '';
    const st = ED_CODE_STATUS.find(s => String(s.ed_code_status_id) === String(edID));
    return st ? String(st.cat_id) : '';
  })();

  row.innerHTML = `
    <form id="${formId}" style="display:none"></form>

    <td>${statusID}</td>
    <td><input form="${formId}" type="text" name="status_name" value="${escapeHtml(statusName||'')}" required></td>

    <td>
      <label>ED_CODE_STATUS_CAT
        <select form="${formId}" id="${edCatSelId}" name="ed_code_status_cat_id">
          <option value="">--</option>
          ${buildOptions(ED_CODE_STATUS_CAT,'ed_code_status_cat_id','descr_en',o=>`${o.descr_en} (ID:${o.ed_code_status_cat_id})`)}
        </select>
      </label>

      <label>ED_CODE_STATUS
        <select form="${formId}" id="${edCodeSelId}" name="ed_code_status_id">
          <option value="">--</option>
        </select>
      </label>
    </td>

    <td>
      <select form="${formId}" id="${gsSelId}" name="gs_code_req_status_id">
        <option value="">--</option>
        ${buildOptions(GS_CODE_REQ_STATUS,'gs_code_req_status_id','descr_en',o=>`${o.descr_en} (ID:${o.gs_code_req_status_id})`)}
      </select>
    </td>

    <td>
      <select form="${formId}" name="is_terminal">
        <option value="0" ${isTerminal? '' : 'selected'}>No</option>
        <option value="1" ${isTerminal? 'selected' : ''}>Yes</option>
      </select>
    </td>

    <td>
      <input form="${formId}" type="number" name="success_path" value="${successPath==null?'':String(successPath)}" placeholder="optional">
    </td>

    <td>
      <button class="icon-btn save" type="button"
              onclick="event.stopPropagation(); saveStatus(${statusID}, ${workflowID}, ${version})"
              title="Save">
        <i class="fa-solid fa-check"></i>
      </button>

      <button class="icon-btn cancel" type="button"
              onclick="event.stopPropagation(); cancelStatusEdit(${statusID}, '${escapeHtml(statusName||'')}', ${edID??'null'}, '${escapeHtml(ED_DESCR_EN||'')}', ${gsID??'null'}, '${escapeHtml(GS_DESCR_EN||'')}', ${isTerminal?1:0}, ${successPath==null?'null':Number(successPath)}, ${workflowID}, ${version})"
              title="Cancel">
        <i class="fa-solid fa-xmark"></i>
      </button>
    </td>
  `;

  const EMPTY_OPTION = { value: '', text: '--' };

  // init Tom Selects
  const tsCat = new TomSelect('#'+edCatSelId, { create:false, closeAfterSelect:true, allowEmptyOption:true });
  const tsEd  = new TomSelect('#'+edCodeSelId,{ 
    create:false, 
    closeAfterSelect:true, 
    allowEmptyOption:true,
    valueField: 'value',
    labelField: 'text',
    searchField: ['text'],
    options: [EMPTY_OPTION]
  });
  const tsGs  = new TomSelect('#'+gsSelId,   { create:false, closeAfterSelect:true, allowEmptyOption:true });

  const catKey = (o) => String(
      o.cat_id ??
      o.ED_CODE_STATUS_CAT_ID ??
      o.ed_code_status_cat_id
    );
  // Populate ED codes filtered by category
  const applyEdCodes = (catVal) => {
       const list = (catVal ? ED_CODE_STATUS.filter(x => catKey(x) === String(catVal)) : [])
          .map(o => ({
            value: String(o.ed_code_status_id),
            text: `${o.descr_en} (ID:${o.ed_code_status_id})`
          }));

        // clear current selection + options, then add fresh list
        tsEd.clear(true);         // clear selection (removes the "--" selection)
        tsEd.clearOptions();      // remove previous options
        tsEd.addOption(EMPTY_OPTION); // add back the empty option
        tsEd.addOptions(list);    // add new options
        tsEd.refreshOptions(false);
        tsEd.setValue('', true);
  };

  tsCat.on('change', (val)=> applyEdCodes(val));
  // initial defaults
  if (edID) {
    tsCat.setValue(String(edCatVal), true);
    applyEdCodes(String(edCatVal));
    tsEd.setValue(String(edID), true);
  }
  if (gsID) tsGs.setValue(String(gsID), true);
}

window.saveStatus = saveStatus;
export function saveStatus(statusID, workflowID, version) {
  const form = document.getElementById('edit_status_form_' + statusID);
  if (!form) return console.error('edit_status_form not found');
  const formData = new FormData(form);

  const edId = formData.get('ed_code_status_id');
  const gsId = formData.get('gs_code_req_status_id');

  if ((!edId && !gsId) || (edId && gsId)) {
    alert("Please select either an ED Code Status or a GS Code Request Status.");
    return;
  }

  const params = new URLSearchParams({
    workflow_id: workflowID,
    version: version,
    status_id: statusID,
    status_name: formData.get('status_name') || '',
    ed_code_status_id: edId || '',
    gs_code_req_status_id: gsId || '',
    is_terminal: formData.get('is_terminal') || '0',
    success_path: formData.get('success_path') || '',
  });

  fetch(`/edit_status?${params.toString()}`, { method: 'GET', headers: { 'Cache-Control': 'no-cache' } })
    .then(r => { if (!r.ok) throw new Error("Failed to update status"); return r.json(); })
    .then(() => {
      const wfid = new URLSearchParams(window.location.search).get('workflow');
      if (wfid) window.viewDetails(wfid);
    })
    .catch(err => console.error("Error updating status:", err));
}

window.cancelStatusEdit = cancelStatusEdit;
export function cancelStatusEdit(statusID, statusName, edID, ED_DESCR_EN, gsID, GS_DESCR_EN, isTerminal, successPath, workflowID, version) {
  const row = document.getElementById(`status_row_${statusID}`);
  if (!row) return;

  row.innerHTML = `
    <td>${statusID}</td>
    <td>${escapeHtml(statusName || '')}</td>
    <td>${edID ? escapeHtml(nullable(ED_DESCR_EN) || '') : ''}</td>
    <td>${gsID ? escapeHtml(nullable(GS_DESCR_EN) || '') : ''}</td>
    <td>${isTerminal ? 'Yes' : 'No'}</td>
    <td>${nullable(successPath)}</td>
    <td>
      <button class="icon-btn edit" type="button"
              onclick="event.stopPropagation(); editStatus(${statusID}, '${escapeHtml(statusName || '')}', ${edID ?? 'null'}, '${escapeHtml(ED_DESCR_EN || '')}', ${gsID ?? 'null'}, '${escapeHtml(GS_DESCR_EN || '')}', ${isTerminal ? 1 : 0}, ${successPath == null ? 'null' : Number(successPath)}, ${workflowID}, ${version})"
              title="Edit">
        <i class="fa-solid fa-pen"></i>
      </button>

      <button class="icon-btn delete" type="button"
              onclick="event.stopPropagation(); deleteStatus(${statusID}, ${workflowID}, ${version})"
              title="Delete">
        <i class="fa-solid fa-trash"></i>
      </button>
    </td>
  `;
}

window.deleteStatus = deleteStatus;
export function deleteStatus(statusID, workflowID, version) {
  if (!confirm("Are you sure you want to delete this status?")) return;
  fetch(`/delete_status?status_id=${statusID}&workflow_id=${workflowID}&version=${version}`)
    .then(r => { if (!r.ok) throw new Error('Failed'); })
    .then(() => {
      const wfid = new URLSearchParams(window.location.search).get('workflow');
      if (wfid) window.viewDetails(wfid);
    })
    .catch(err => console.error("Error deleting status:", err));
}

// ---------- helpers ----------
window.getCat = getCat;
export function getCat(edCodeStatusId) {
  const st = ED_CODE_STATUS.find(s => Number(s.ed_code_status_id) === Number(edCodeStatusId));
  if (!st) return '';
  const cat = ED_CODE_STATUS_CAT.find(c => Number(c.ed_code_status_cat_id) === Number(st.cat_id));
  return cat ? `${cat.ed_code_status_cat_id}_${cat.descr_en}` : '';
}

// legacy fallback (kept)
window.isEdinCat = isEdinCat;
export function isEdinCat(edStatus) {
  let input = document.getElementById('ed_code_status_cat_edit_input');
  if (!input) input = document.getElementById('ed_code_status_cat_input');
  if (!input || !input.value) return true;
  const catId = Number(input.value.split('_')[0]);
  return Number(edStatus.cat_id) === catId;
}

// per-row versions used above
export function isEdinCatByInputId(edStatus, inputId) {
  const input = document.getElementById(inputId);
  if (!input || !input.value) return true;
  const catId = Number(input.value.split('_')[0]);
  return Number(edStatus.cat_id) === catId;
}

window.updateEdCodes = updateEdCodes; // keep existing name
export function updateEdCodes(edCodesDatalist) {
  // kept for backwards compat (expects an element)
  if (!edCodesDatalist) return;
  const input = document.getElementById('ed_code_status_cat_edit_input');
  edCodesDatalist.innerHTML = ED_CODE_STATUS
    .filter(o => !input || !input.value || Number(input.value.split('_')[0]) === Number(o.cat_id))
    .map(o => `<option value="${o.ed_code_status_id}_${escapeHtml(o.descr_en)}">${o.ed_code_status_id}_${escapeHtml(o.descr_en)}</option>`)
    .join('');
}

// new, safer API used in the template above
window.updateEdCodesByIds = updateEdCodesByIds;
export function updateEdCodesByIds(catInputId, listId) {
  const list = document.getElementById(listId);
  if (!list) return;
  list.innerHTML = ED_CODE_STATUS
    .filter(o => isEdinCatByInputId(o, catInputId))
    .map(o => `<option value="${o.ed_code_status_id}_${escapeHtml(o.descr_en)}">${o.ed_code_status_id}_${escapeHtml(o.descr_en)}</option>`)
    .join('');
}
