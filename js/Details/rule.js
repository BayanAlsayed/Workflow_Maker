import { WF_STATUSES, SE_CODE_USER_TYPE, SE_ACCNT, setUserType } from "./details.js";
import { escapeHtml, nullable } from "../helpers.js";

window.createRule = createRule;
export function createRule(e, workflowID) {
  if (e && typeof e.preventDefault === 'function') e.preventDefault();

  const form = document.getElementById('add-rule-form');
  const formData = new FormData(form);
  const data = Object.fromEntries(formData.entries());
  data.workflow_id = workflowID;

  fetch('/add_rule', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  })
  .then(response => {
    if (response.ok) {
      console.log("Rule created");
      viewDetails(new URLSearchParams(window.location.search).get('workflow'));
    } else {
      console.error("Failed to create rule");
    }
  })
  .catch(error => {
    console.error("Error creating rule:", error);
  });
}

window.editRule = editRule;
export function editRule(rule_id, from_status_id, from_status_name, to_status_id, to_status_name, se_code_user_type_id, user_type_en, se_accnt_id, accnt_en, action_button, action_function, workflowID) {
    console.log("edit status for workflow ID:", workflowID);

    const row = document.getElementById(`rule_row_${rule_id}`);
    if (!row) return;
  
    // per-row unique ids (avoid clashes if multiple edits)
    const fromStatusListId  = `from_status_edit_${rule_id}`;
    const toStatusListId    = `to_status_edit_${rule_id}`;

    const fromStatusValue = from_status_id ? `${from_status_id}_${from_status_name || ''}` : '';
    const toStatusValue = to_status_id ? `${to_status_id}_${to_status_name || ''}` : '';

    const formId = `edit_rule_form_${rule_id}`;
    row.innerHTML = `
      <!-- hidden form element to own the inputs -->
      <form id="${formId}" style="display:none"></form>

      <td>${rule_id}</td>
      <td>
        <input form="${formId}" list="${fromStatusListId}" 
        name="from_status_id" placeholder="Search..."
        value="${escapeHtml(fromStatusValue || '')}" required>
        <datalist id="${fromStatusListId}">
          ${WF_STATUSES.map(s =>
            `<option value="${s.status_id}_${escapeHtml(s.status_name)}">${s.status_id}_${escapeHtml(s.status_name)}</option>`
          ).join('')}
        </datalist>
      </td>

      <td>
        <input form="${formId}" list="${toStatusListId}" 
        name="to_status_id" placeholder="Search..."
        value="${escapeHtml(toStatusValue || '')}" required>
        <datalist id="${toStatusListId}">
          ${WF_STATUSES.map(s =>
            `<option value="${s.status_id}_${escapeHtml(s.status_name)}">${s.status_id}_${escapeHtml(s.status_name)}</option>`
          ).join('')}
        </datalist>
      </td>

      <td>
        <input form="${formId}" list="user_types" name="se_code_user_type"
        id="se_code_user_type_edit_${rule_id}"
        placeholder="Search..." 
        value="${se_code_user_type_id ? `${se_code_user_type_id}_${escapeHtml(user_type_en || '')}` : ''}"
        required>
        <datalist id="user_types">
          ${SE_CODE_USER_TYPE.map(ut =>
            `<option value="${ut.se_code_user_type_id}_${escapeHtml(ut.descr_en)}">${ut.se_code_user_type_id}_${escapeHtml(ut.descr_en)}</option>`
          ).join('')}
        </datalist>
      </td>

      <td>
        <input form="${formId}" list="accnts" name="se_accnt" 
        placeholder="Search..." 
        value="${se_accnt_id ? `${se_accnt_id}_${escapeHtml(accnt_en || '')}` : ''}"
        onchange="setUserType(this, 'se_code_user_type_edit_${rule_id}')" >
        <datalist id="accnts">
          ${SE_ACCNT.map(acct =>
            `<option value="${acct.se_accnt_id}_${escapeHtml(acct.descr_en)}">${acct.se_accnt_id}_${escapeHtml(acct.descr_en)}</option>`
          ).join('')}
        </datalist>
      </td>

      <td>
        <input form="${formId}" type="text" name="action_button" 
        value="${escapeHtml(action_button || '')}" required>
      </td>

      <td>
        <input form="${formId}" type="text" name="action_function" 
        value="${escapeHtml(action_function || '')}">
      </td>

      <td>
        <button class="icon-btn save" type="button"
                onclick="event.stopPropagation(); saveRule(${rule_id}, ${workflowID})"
                title="Save">
          <i class="fa-solid fa-check"></i>
        </button>
        <button class="icon-btn cancel" type="button"
                onclick="event.stopPropagation(); cancelRuleEdit(${rule_id}, '${from_status_id}', '${escapeHtml(from_status_name || '')}', '${to_status_id}', '${escapeHtml(to_status_name || '')}', '${se_code_user_type_id}', '${escapeHtml(user_type_en || '')}', '${se_accnt_id}', '${escapeHtml(accnt_en || '')}', '${escapeHtml(action_button || '')}', '${escapeHtml(action_function || '')}', ${workflowID})"
                title="Cancel">
          <i class="fa-solid fa-xmark"></i>
        </button>
      </td>

    `;
}

window.deleteRule = deleteRule;
export function deleteRule(rule_id, workflowID) {
  if (!confirm("Are you sure you want to delete this rule?")) return;

  fetch(`/delete_rule?rule_id=${rule_id}&workflow_id=${workflowID}`)
    .then(response => {
      if (response.ok) {
        console.log("Rule deleted:", rule_id);
        // refresh only the workflow details, not full reload
        const wfid = new URLSearchParams(window.location.search).get('workflow');
        if (wfid) window.viewDetails(wfid);
      } else {
        console.error("Failed to delete rule:", rule_id);
      }
    })
    .catch(error => {
      console.error("Error deleting rule:", error);
    });
}

window.saveRule = saveRule;
export function saveRule(rule_id, workflowID){
  console.log("save rule for workflow ID:", workflowID);

  const form = document.getElementById('edit_rule_form_' + rule_id);
  if (!form) return console.error('edit_rule_form not found');

  const formData = new FormData(form);
  console.log('formData:', Object.fromEntries(formData.entries()));

  const idFromCombo = (v) => {
    if (!v) return '';
    const [id] = String(v).split('_', 1);
    return id;
  };

  const params = new URLSearchParams({
    workflow_id: workflowID,
    rule_id: rule_id,
    from_status_id: idFromCombo(formData.get('from_status_id')),
    to_status_id: idFromCombo(formData.get('to_status_id')),
    se_code_user_type_id: idFromCombo(formData.get('se_code_user_type')),
    se_accnt_id: idFromCombo(formData.get('se_accnt')),
    action_button: formData.get('action_button'),
    action_function: nullable(formData.get('action_function')),
  });

  console.log('Saving rule with params:', params.toString());

  fetch(`/edit_rule?${params.toString()}`, {
    method: 'GET',
    headers: { 'Cache-Control': 'no-cache' },
  })
    .then(response => {
      if (!response.ok) throw new Error("Failed to update rule");
      return response.json();
    })
    .then(() => {
      console.log("Rule updated:", rule_id);
      // refresh only the workflow details, not full reload
      const wfid = new URLSearchParams(window.location.search).get('workflow');
      if (wfid) window.viewDetails(wfid);
    })
    .catch(error => console.error("Error updating rule:", error));
}



window.cancelRuleEdit = cancelRuleEdit;
export function cancelRuleEdit(rule_id, from_status_id, from_status_name, to_status_id, to_status_name, se_code_user_type_id, user_type_en, se_accnt_id, accnt_en, action_button, action_function, workflowID){
  console.log("cancel rule edit for workflow ID:", workflowID);

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
        onclick="event.stopPropagation(); editRule(${rule_id}, ${from_status_id}, '${from_status_name}', ${to_status_id}, '${to_status_name}', ${se_code_user_type_id}, '${user_type_en}', ${se_accnt_id}, '${accnt_en}', '${action_button}', '${action_function}', ${workflowID})"
        title="Edit">
        <i class="fa-solid fa-pen"></i>
      </button>

      <button class="icon-btn delete" type="button"
        onclick="event.stopPropagation(); deleteRule(${rule_id}, ${workflowID})"
        title="Delete">
        <i class="fa-solid fa-trash"></i>
      </button>
    </td>
  `;
}


