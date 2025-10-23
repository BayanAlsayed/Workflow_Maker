import { WF_STATUSES } from "./details.js";
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
    const edCatListId       = `ed_codes_cat_edit_${rule_id}`;
    const edCodesListId     = `ed_codes_edit_${rule_id}`;
    const gsCodesListId     = `gs_codes_${rule_id}`;

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
        <input form="${formId}" list="user_types" name="se_code_user_type_id" 
        placeholder="Search..." 
        value="${se_code_user_type_id ? `${se_code_user_type_id}_${escapeHtml(user_type_en || '')}` : ''}">
      </td>

      <td>
        <input form="${formId}" list="accnts" name="se_accnt_id" 
        placeholder="Search..." 
        value="${se_accnt_id ? `${se_accnt_id}_${escapeHtml(accnt_en || '')}` : ''}">
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
                onclick="event.stopPropagation(); cancelRuleEdit(${rule_id}, '${escapeHtml(statusName || '')}', ${edID ?? 'null'}, '${escapeHtml(ED_DESCR_EN || '')}', ${gsID ?? 'null'}, '${escapeHtml(GS_DESCR_EN || '')}', ${isTerminal ? 1 : 0}, ${successPath == null ? 'null' : Number(successPath)}, ${workflowID})"
                title="Cancel">
          <i class="fa-solid fa-xmark"></i>
        </button>
      </td>

    `;
}

window.deleteRule = deleteRule;
export function deleteRule(){

}

window.saveRule = saveRule;
export function saveRule(){

}

window.cancelRuleEdit = cancelRuleEdit;
export function cancelRuleEdit(){

}


