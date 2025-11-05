import { escapeHtml, nullable, opt, buildOptions, initTomSelectSingle } from "../helpers.js";
import { WF_CONDITIONS } from "./details.js";


window.viewConditions = viewConditions;
export function viewConditions() {
  const container = document.getElementById('main_container');
  container.innerHTML = ''; // Clear previous content

  fetch('/view_conditions', { method: 'GET', headers: { 'Cache-Control': 'no-cache' } })
    .then(response => response.json())
    .then(data => {
        if(data) {
            const conditionsTable= `
                <table id="conditions_table" class="wf-table">
                    <thead>
                        <tr>
                            <th>Function Name</th>
                            <th>Description</th>
                            <th>Active Workflows Using</th>
                            <th>Active Rules Count</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${data.map(cond => `
                            <tr id="condition_row_${cond.wf_condition_id}">
                                <td>${cond.func_name}</td>
                                <td>${cond.description}</td>
                                <td>${cond.active_workflows_using}</td>
                                <td>${cond.active_rules_count}</td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            `;
            container.innerHTML=conditionsTable;
        } else {
            container.innerHTML=`there are no conditions yet!`
        }

        const conditionsForm = `
            <div class="wf-create">
                <form id="add_condition_form" onsubmit="createCondition(event)">
                    <label for="function_name">Function Name:</label>
                    <input type="text" id="function_name" name="function_name" required>
                    <label for="description">Description:</label>
                    <textarea id="description" name="description" required></textarea>
                    <button type="submit">Add Condition</button>
                </form>
            </div>
        `;

        container.innerHTML += conditionsForm;

    });

  container.style.display = 'block'; // Show the container
}

window.createCondition = createCondition;
export function createCondition(e) {
    e.preventDefault();

    const form = document.getElementById("add_condition_form");
    const formData = new FormData(form);
    const data = Object.fromEntries(formData.entries());

    fetch('/add_condition', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    })
    .then(r => { if (!r.ok) throw new Error('Failed to create condition'); })
    .then(() => {
        viewConditions();
    })
    .catch(err => console.error('Error creating condition:', err));

}

window.addConditionPrompt = addConditionPrompt;
export function addConditionPrompt(rule_id, workflowID, version) {
    const ruleRow = document.getElementById(`rule_row_${rule_id}`);
    if (!ruleRow) return;

    const formId = `add_rule_condition_form_${rule_id}`;
    const con_sel_id = `condition_select_${rule_id}`;
    const desc_id = `condition_description_${rule_id}`;
    const type_sel_id = `type_select_${rule_id}`;

    ruleRow.insertAdjacentHTML('afterend', `
        <tr id="condition_form_row_${rule_id}">
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
                <button type="button" onclick="cancelConditionPrompt(${rule_id})">Cancel</button>
            </td>
        </tr>
    `);

    const tsCon = initTomSelectSingle(`#${con_sel_id}`);

    tsCon.on('change', (val) => {
        updateCondDesc(val, desc_id);
    });
}

window.updateCondDesc = updateCondDesc;
export function updateCondDesc(conID, descID) {
    const descEl = document.getElementById(descID);
    const condition = WF_CONDITIONS.find(c => c.wf_condition_id == conID);
    if (condition) {
        descEl.value = condition.description;
    } else {
        descEl.value = 'Condition description will be shown here';
    }
}

window.cancelConditionPrompt = cancelConditionPrompt;
export function cancelConditionPrompt(rule_id) {
    const formRow = document.getElementById(`condition_form_row_${rule_id}`);
    if (formRow) {
        formRow.remove();
    }
}

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
