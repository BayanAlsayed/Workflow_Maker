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



