import { escapeHtml, nullable, opt, buildOptions, initTomSelectSingle } from "../helpers.js";
import { WF_CONDITIONS } from "./details.js";

window.addRuleCondition = addRuleCondition;
export function addRuleCondition(rule_id, workflowID, version) {
    const formId = `add_rule_condition_form_${rule_id}`;
    const form = document.getElementById(formId);
    const formData = new FormData(form);
    const data = Object.fromEntries(formData.entries());
    data.rule_id = rule_id;
    data.workflow_id = workflowID;
    data.version = version;

    fetch('/add_rule_condition', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    })
    .then(r => { if (!r.ok) throw new Error('Failed to add condition to rule'); })
    .then(() => {
        viewDetails(workflowID, version);
    })
    .catch(err => console.error('Error adding condition to rule:', err));
}

window.toggleRuleConditions = toggleRuleConditions;
export function toggleRuleConditions(ruleID) {
  const wrapper = document.getElementById(`rule_conditions_wrapper_${ruleID}`);
  const plusIcon = document.getElementById(`icon_plus_${ruleID}`);
  if (!wrapper) return;
  if (wrapper.style.display == 'none') {
    wrapper.style.display = 'table-row'
    plusIcon.style.transform = 'rotate(45deg)'
  } else {
    wrapper.style.display = 'none';
    plusIcon.style.transform = 'rotate(0deg)'
  }
};

window.addRuleConditionPrompt = addRuleConditionPrompt;
export function addRuleConditionPrompt(rule_id, workflowID, version) {
    const ruleConditionRow = document.getElementById(`rule_conditions_wrapper_${rule_id}`);
    if (!ruleConditionRow) return;

    const ruleConditionForm = document.getElementById(`rule_condition_form_row_${rule_id}`);
    if (ruleConditionForm) return;

    const formId = `add_rule_condition_form_${rule_id}`;
    const con_sel_id = `condition_select_${rule_id}`;
    const desc_id = `condition_description_${rule_id}`;
    const type_sel_id = `type_select_${rule_id}`;

    ruleConditionRow.insertAdjacentHTML('afterend', `
        <tr id="rule_condition_form_row_${rule_id}">
            <form id="${formId}" style="display:none"></form>

            <td colspan="2">
                    <label for="${con_sel_id}">Select Condition:</label>
                    <select form="${formId}" id="${con_sel_id}" name="wf_condition_id" required>
                        <option value="">--Select Condition--</option>
                        ${buildOptions(WF_CONDITIONS, 'wf_condition_id', 'func_name', o => `${o.func_name}`)}
                    </select>
                    
            </td>
            <td colspan="4">
                <label for="${desc_id}">Condition Description:</label>
                <textarea id="${desc_id}" name="description" readonly>Condition description will be shown here</textarea>
            </td>
            <td>
                <label for="${type_sel_id}">Condition Type:</label>
                <select form="${formId}" id="${type_sel_id}" name="condition_type" required>
                    <option value="pre">pre</option>
                    <option value="post">post</option>
                </select>
            </td>
            <td>
                <button type="button" onclick="addRuleCondition(${rule_id}, ${workflowID}, ${version})">Add Condition to Rule</button>
                <button type="button" onclick="cancelRuleConditionPrompt(${rule_id})">Cancel</button>
            </td>
        </tr>
    `);

    const tsCon = initTomSelectSingle(`#${con_sel_id}`);

    tsCon.on('change', (val) => {
        updateCondDesc(val, desc_id);
    });
}

window.cancelRuleConditionPrompt = cancelRuleConditionPrompt;
export function cancelRuleConditionPrompt(rule_id) {
    const formRow = document.getElementById(`rule_condition_form_row_${rule_id}`);
    if (formRow) {
        formRow.remove();
    }
}

window.deleteRuleCondition = deleteRuleCondition;
export function deleteRuleCondition(wf_rule_condition_id, workflowID, version) {
    if (!confirm('Are you sure you want to delete this rule condition?')) return;

    fetch(`/delete_rule_condition?id=${wf_rule_condition_id}`)
    .then(r => { if (!r.ok) throw new Error('Failed to delete rule condition'); })
    .then(() => {
        viewDetails(workflowID, version);
    })
    .catch(err => console.error('Error deleting rule condition:', err));
}

window.editRuleCondition = editRuleCondition;
export function editRuleCondition(wf_rule_condition_id, rule_condition_id, rule_id, wf_condition_id, func_name, description, type, workflowID, version) {
    const conditionCard = document.getElementById(`cond_card_${rule_condition_id}_${rule_id}`);
    if (!conditionCard) return;

    const formId = `edit_rule_condition_form_${rule_condition_id}`;
    const con_sel_id = `edit_condition_select_${rule_condition_id}`;
    const desc_id = `edit_condition_description_${rule_condition_id}`;
    const type_sel_id = `edit_type_select_${rule_condition_id}`;
    conditionCard.innerHTML= `
        <form id="${formId}"></form>
        <div class="cond-func">
            <label for="${con_sel_id}">Select Condition:</label>
            <select form="${formId}" id="${con_sel_id}" name="wf_condition_id" required>
                <option value="">--Select Condition--</option>
                ${buildOptions(WF_CONDITIONS, 'wf_condition_id', 'func_name', o => `${o.func_name}`)}
            </select>
        </div>
        <div class="cond-desc">
            <label for="${desc_id}">Condition Description:</label>
            <textarea id="${desc_id}" name="description" readonly>${description}</textarea>
        </div>
        <div class="cond-meta">
            <div class="cond-left">
            <span>
                <label for="${type_sel_id}">Condition Type:</label>
                <select form="${formId}" id="${type_sel_id}" name="condition_type" required>
                    <option value="pre" ${type === 'pre' ? 'selected' : ''}>pre</option>
                    <option value="post" ${type === 'post' ? 'selected' : ''}>post</option>
                </select>
            </span>
            </div>
            <div class="cond-actions">
            <button class="icon-btn cond" title="save rule condition" onclick="event.stopPropagation(); saveRuleCondition(${wf_rule_condition_id}, ${rule_condition_id}, ${workflowID}, ${version})"><i class="fa-solid fa-check"></i></button>
            <button class="icon-btn cond" title="cancel edit" onclick="event.stopPropagation(); cancelRuleConditionEdit(${wf_rule_condition_id}, ${rule_condition_id}, ${rule_id}, ${wf_condition_id}, '${func_name}', '${description}', '${type}', ${workflowID}, ${version})"><i class="fa-solid fa-xmark"></i></button>
            </div>
        </div>
    `
    const tsCon = initTomSelectSingle(`#${con_sel_id}`);
    if (wf_condition_id) {
        tsCon.setValue(String(wf_condition_id), true);
    }

    tsCon.on('change', (val) => {
        updateCondDesc(val, desc_id);
    });
}

window.cancelRuleConditionEdit = cancelRuleConditionEdit;
export function cancelRuleConditionEdit(wf_rule_condition_id, rule_condition_id, rule_id, wf_condition_id, func_name, description, type, workflowID, version) {
    const conditionCard = document.getElementById(`cond_card_${rule_condition_id}_${rule_id}`);
    if (!conditionCard) return;

    conditionCard.innerHTML = `
        <div class="cond-func">${escapeHtml(func_name || '(Unnamed)')}</div>
        <div class="cond-desc">${escapeHtml(description || '')}</div>
        <div class="cond-meta">
            <div class="cond-left">
            <span>${escapeHtml(type || '')}</span>
            </div>
            <div class="cond-actions">
            <button class="icon-btn cond" title="Edit condition" onclick="event.stopPropagation(); editRuleCondition(${wf_rule_condition_id}, ${rule_condition_id}, ${rule_id}, ${wf_condition_id}, '${func_name}', '${description}', '${type}', ${workflowID}, ${version})"><i class="fa-solid fa-pen"></i></button>
            <button class="icon-btn cond" title="Delete condition" onclick="event.stopPropagation(); deleteRuleCondition(${wf_rule_condition_id}, ${workflowID}, ${version})"><i class="fa-solid fa-trash"></i></button>
            </div>
        </div>
    `;
}

window.saveRuleCondition = saveRuleCondition;
export function saveRuleCondition(wf_rule_condition_id, rule_condition_id, workflow_id, version) {
    const formId = `edit_rule_condition_form_${rule_condition_id}`;
    const form = document.getElementById(formId);
    if (!form ) return console.error('edit_rule_condition_form is not found');
    const formData = new FormData(form);

    const params = new URLSearchParams({
        wf_rule_condition_id: wf_rule_condition_id,
        wf_condition_id: formData.get('wf_condition_id'),
        condition_type: formData.get('condition_type'),
    });
    

    fetch(`/edit_rule_condition?${params.toString()}`)
    .then(r => { if (!r.ok) throw new Error('Failed to save rule condition'); })
    .then(() => {
        viewDetails(workflow_id, version);
    })
    .catch(err => console.error('Error saving rule condition:', err));
}

