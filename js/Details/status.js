import { escapeHtml, nullable } from '../helpers.js';
import { ED_CODE_STATUS_CAT, ED_CODE_STATUS, GS_CODE_REQ_STATUS } from './details.js';

window.createStatus = createStatus;
export function createStatus(e, workflowID) {
  if (e && typeof e.preventDefault === 'function') e.preventDefault();

  const form = document.getElementById('add-status-form');
  if (!form) return console.error('add-status-form not found');

  const formData = new FormData(form);
  const edCodeStatusId = formData.get('ed_code_status_id');
  const gsCodeReqStatusId = formData.get('gs_code_req_status_id');

  if ((!edCodeStatusId && !gsCodeReqStatusId) || (edCodeStatusId && gsCodeReqStatusId)) {
    alert("Please select either an ED Code Status or a GS Code Request Status.");
    return;
  }

  const data = Object.fromEntries(formData.entries());
  data.workflow_id = workflowID;

  fetch('/add_status', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  .then(r => {
    if (!r.ok) throw new Error('Failed to create status');
    const wfid = new URLSearchParams(window.location.search).get('workflow');
    if (wfid) window.viewDetails(wfid);
  })
  .catch(err => console.error('Error creating status:', err));
}

window.editStatus = editStatus;
export function editStatus(statusID, statusName, edID, ED_DESCR_EN, gsID, GS_DESCR_EN, isTerminal, successPath, workflowID) {
  console.log("edit status for workflow ID:", workflowID);

  const row = document.getElementById(`status_row_${statusID}`);
  if (!row) return;

  // per-row unique ids (avoid clashes if multiple edits)
  const edCatInputId  = `ed_code_status_cat_edit_input_${statusID}`;
  const edCatListId   = `ed_codes_cat_edit_${statusID}`;
  const edCodesListId = `ed_codes_edit_${statusID}`;
  const gsCodesListId = `gs_codes_${statusID}`;

  const edCatValue = edID ? getCat(edID) : '';
  const edCodeValue = edID ? `${edID}_${ED_DESCR_EN || ''}` : '';
  const gsCodeValue = gsID ? `${gsID}_${GS_DESCR_EN || ''}` : '';

  const formId = `edit_status_form_${statusID}`;
  row.innerHTML = `
    <!-- hidden form element to own the inputs -->
    <form id="${formId}" style="display:none"></form>

    <td>${statusID}</td>
    <td><input form="${formId}" type="text" name="status_name"
              value="${escapeHtml(statusName || '')}" required></td>

    <td>
      <label>ED_CODE_STATUS_CAT
        <input form="${formId}" list="${edCatListId}" id="${edCatInputId}"
              name="ed_code_status_cat_id" placeholder="Search..."
              value="${escapeHtml(edCatValue)}"
              onchange="updateEdCodesByIds('${edCatInputId}', '${edCodesListId}')">
        <datalist id="${edCatListId}">
          ${ED_CODE_STATUS_CAT.map(o =>
            `<option value="${o.ed_code_status_cat_id}_${escapeHtml(o.descr_en)}">${o.ed_code_status_cat_id}_${escapeHtml(o.descr_en)}</option>`
          ).join('')}
        </datalist>
      </label>

      <label>ED_CODE_STATUS
        <input form="${formId}" list="${edCodesListId}" name="ed_code_status_id"
              placeholder="Search..." value="${escapeHtml(edCodeValue)}">
        <datalist id="${edCodesListId}">
          ${ED_CODE_STATUS.filter(o => isEdinCatByInputId(o, edCatInputId)).map(o =>
            `<option value="${o.ed_code_status_id}_${escapeHtml(o.descr_en)}">${o.ed_code_status_id}_${escapeHtml(o.descr_en)}</option>`
          ).join('')}
        </datalist>
      </label>
    </td>

    <td>
      <input form="${formId}" list="${gsCodesListId}" name="gs_code_req_status_id"
            placeholder="Search..." value="${escapeHtml(gsCodeValue)}">
      <datalist id="${gsCodesListId}">
        ${GS_CODE_REQ_STATUS.map(o =>
          `<option value="${o.gs_code_req_status_id}_${escapeHtml(o.descr_en)}">${o.gs_code_req_status_id}_${escapeHtml(o.descr_en)}</option>`
        ).join('')}
      </datalist>
    </td>

    <td>
      <select form="${formId}" name="is_terminal">
        <option value="0" ${isTerminal ? '' : 'selected'}>No</option>
        <option value="1" ${isTerminal ? 'selected' : ''}>Yes</option>
      </select>
    </td>

    <td>
      <input form="${formId}" type="number" name="success_path"
            placeholder="optional" value="${successPath == null ? '' : String(successPath)}">
    </td>

    <td>
      <button class="icon-btn save" type="button"
              onclick="event.stopPropagation(); saveStatus(${statusID}, ${workflowID})"
              title="Save">
        <i class="fa-solid fa-check"></i>
      </button>
      <button class="icon-btn cancel" type="button"
              onclick="event.stopPropagation(); cancelStatusEdit(${statusID}, '${escapeHtml(statusName || '')}', ${edID ?? 'null'}, '${escapeHtml(ED_DESCR_EN || '')}', ${gsID ?? 'null'}, '${escapeHtml(GS_DESCR_EN || '')}', ${isTerminal ? 1 : 0}, ${successPath == null ? 'null' : Number(successPath)}, ${workflowID})"
              title="Cancel">
        <i class="fa-solid fa-xmark"></i>
      </button>
    </td>
  `;

}

window.saveStatus = saveStatus;
export function saveStatus(statusID, workflowID) {
  console.log("save status for workflow ID:", workflowID);

  const form = document.getElementById('edit_status_form_' + statusID);
  if (!form) return console.error('edit_status_form not found');

  const formData = new FormData(form);
  const edRaw = formData.get('ed_code_status_id');
  const gsRaw = formData.get('gs_code_req_status_id');

  console.log('formData:', Object.fromEntries(formData.entries()));


  if ((!edRaw && !gsRaw) || (edRaw && gsRaw)) {
    alert("Please select either an ED Code Status or a GS Code Request Status.");
    return;
  }

  const idFromCombo = (v) => {
    if (!v) return '';
    const [id] = String(v).split('_', 1);
    return id;
  };

  // Build query parameters for GET
  const params = new URLSearchParams({
    workflow_id: workflowID,
    status_id: statusID,
    status_name: formData.get('status_name') || '',
    ed_code_status_id: idFromCombo(edRaw),
    gs_code_req_status_id: idFromCombo(gsRaw),
    is_terminal: formData.get('is_terminal') || '0',
    success_path: formData.get('success_path') || '',
  });

  fetch(`/edit_status?${params.toString()}`, {
    method: 'GET',
    headers: { 'Cache-Control': 'no-cache' },
  })
    .then(response => {
      if (!response.ok) throw new Error("Failed to update status");
      return response.json();
    })
    .then(() => {
      console.log("Status updated:", statusID);
      // refresh only the workflow details, not full reload
      const wfid = new URLSearchParams(window.location.search).get('workflow');
      if (wfid) window.viewDetails(wfid);
    })
    .catch(error => console.error("Error updating status:", error));
}


window.cancelStatusEdit = cancelStatusEdit;
export function cancelStatusEdit(statusID, statusName, edID, ED_DESCR_EN, gsID, GS_DESCR_EN, isTerminal, successPath, workflowID) {
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
              onclick="event.stopPropagation(); editStatus(${statusID}, '${escapeHtml(statusName || '')}', ${edID ?? 'null'}, '${escapeHtml(ED_DESCR_EN || '')}', ${gsID ?? 'null'}, '${escapeHtml(GS_DESCR_EN || '')}', ${isTerminal ? 1 : 0}, ${successPath == null ? 'null' : Number(successPath)}, ${workflowID})"
              title="Edit">
        <i class="fa-solid fa-pen"></i>
      </button>

      <button class="icon-btn delete" type="button"
              onclick="event.stopPropagation(); deleteStatus(${statusID}, ${workflowID})"
              title="Delete">
        <i class="fa-solid fa-trash"></i>
      </button>
    </td>
  `;
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
