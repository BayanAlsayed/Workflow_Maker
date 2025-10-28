import { escapeHtml, nullable } from "../helpers.js";
export var ED_CODE_STATUS_CAT = [];
export var ED_CODE_STATUS = [];
export var GS_CODE_REQ_STATUS = [];
export var SE_CODE_USER_TYPE = [];
export var SE_ACCNT = [];
export var WF_STATUSES = [];

window.viewDetails = viewDetails;
export function viewDetails(workflowID) {
  history.pushState({}, '', `?workflow=${workflowID}`);

  console.log("Viewing details for workflow ID:", workflowID);
  const container = document.getElementById('rules_container');
  container.innerHTML = ''; // Clear previous content
  container.style.display = 'block';

  fetch(`/view_workflow/${workflowID}`)
    .then(response => {
      if (!response.ok) throw new Error("Failed to fetch workflow details");
      return response.json();
    })
    .then(data => {
      console.log("Workflow details:", data);

      const statusesTable = `
        <h3>Statuses</h3>
        <table class="wf-table">
          <thead>
            <tr>
              <th>STATUS_ID</th>
              <th>STATUS_NAME</th>
              <th>ED_DESCR_EN</th>
              <th>GS_DESCR_EN</th>
              <th>IS_TERMINAL</th>
              <th>SUCCESS_PATH</th>
              <th>ACTIONS</th>
            </tr>
          </thead>
          <tbody>
            ${!data.statuses ? `
              <tr><td colspan="9">No statuses found.</td></tr>
            ` : data.statuses.map(s => `
              <tr id="status_row_${s.status_id}">
                <td>${s.status_id}</td>
                <td>${escapeHtml(s.status_name)}</td>
                <td>${nullable(s.ed_descr_en)}</td>
                <td>${nullable(s.gs_descr_en)}</td>
                <td>${s.is_terminal ? 'Yes' : 'No'}</td>
                <td>${nullable(s.success_path)}</td>
                <td>
                  <button class="icon-btn edit" type="button"
                    onclick="event.stopPropagation(); editStatus(${s.status_id}, '${s.status_name}', ${s.ed_code_status_id}, '${s.ed_descr_en}', ${s.gs_code_req_status_id}, '${s.gs_descr_en}', ${s.is_terminal}, ${s.success_path}, ${workflowID})"
                    title="Edit">
                    <i class="fa-solid fa-pen"></i>
                  </button>

                  <button class="icon-btn delete" type="button"
                    onclick="event.stopPropagation(); deleteStatus(${s.status_id}, ${workflowID})"
                    title="Delete">
                    <i class="fa-solid fa-trash"></i>
                  </button>
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      `;

      const rulesTable = `
        <h3>Rules</h3>
        <table class="wf-table">
          <thead>
            <tr>
              <th>RULE_ID</th>
              <th>FROM_STATUS</th>
              <th>TO_STATUS</th>
              <th>USER_TYPE</th>
              <th>ACCOUNT</th>
              <th>ACTION_BUTTON</th>
              <th>ACTION_FUNCTION</th>
              <th>ACTIONS</th>
            </tr>
          </thead>
          <tbody>
            ${!data.rules ? `
              <tr><td colspan="8">No rules found.</td></tr>
            ` : data.rules.map(r => `
              <tr id="rule_row_${r.rule_id}">
                <td>${r.rule_id}</td>
                <td>${nullable(r.from_status_name)}</td>
                <td>${nullable(r.to_status_name)}</td>
                <td>${nullable(r.user_type_en)}</td>
                <td>${nullable(r.accnt_en)}</td>
                <td>${escapeHtml(r.action_button || '')}</td>
                <td>${escapeHtml(r.action_function || '')}</td>
                <td>
                  <button class="icon-btn edit" type="button"
                    onclick="event.stopPropagation(); editRule(${r.rule_id}, ${r.from_status_id}, '${r.from_status_name}', ${r.to_status_id}, '${r.to_status_name}', ${r.se_code_user_type_id}, '${r.user_type_en}', ${r.se_accnt_id}, '${r.accnt_en}', '${r.action_button}', '${r.action_function}', ${workflowID})"
                    title="Edit">
                    <i class="fa-solid fa-pen"></i>
                  </button>

                  <button class="icon-btn delete" type="button"
                    onclick="event.stopPropagation(); deleteRule(${r.rule_id}, ${workflowID})"
                    title="Delete">
                    <i class="fa-solid fa-trash"></i>
                  </button>
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      `;

      const statusFormsHTML = `<div class="wf-create">
          <h3>Add Status</h3>
          <form id="add-status-form" onsubmit="createStatus(event, ${workflowID})">
            <div class="grid">
              <label>STATUS_NAME
                <input type="text" name="status_name" required />
              </label>
              <label>ED_CODE_STATUS_CAT
                <input list="ed_codes_cat" name="ed_code_status_cat_id" id="ed_code_status_cat_input" onchange="updateEdCodesByIds('ed_code_status_cat_input', 'ed_codes')" placeholder="Search..." />
                <datalist id="ed_codes_cat"></datalist>
              </label>
              <label>ED_CODE_STATUS
                <input list="ed_codes" name="ed_code_status_id" placeholder="Search..." />
                <datalist id="ed_codes"></datalist>
              </label>
              <label>GS_REQ_STATUS
                <input list="gs_codes" name="gs_code_req_status_id" placeholder="Search..." />
                <datalist id="gs_codes"></datalist>
              </label>
              <label>IS_TERMINAL
                <select name="is_terminal">
                  <option value="0">No</option>
                  <option value="1">Yes</option>
                </select>
              </label>
              <label>SUCCESS_PATH
                <input type="number" name="success_path" placeholder="optional" />
              </label>
            </div>
            <button type="submit" class="btn primary">Create Status</button>
          </form>
        </div>`;
      // After you build statusesTable and rulesTable:
      const ruleFormsHTML = `
        <div class="wf-create">
          <h3>Add Rule</h3>
          <form id="add-rule-form" onsubmit="createRule(event, ${workflowID})">
            <div class="grid">
              <label>FROM STATUS
                <input list="wf_statuses" name="from_status" required placeholder="Search..." />
              </label>
              <label>TO STATUS
                <input list="wf_statuses" name="to_status" required placeholder="Search..." />
              </label>
              <datalist id="wf_statuses"></datalist>

              <label>USER TYPE
                <input list="user_types" name="se_code_user_type" id="se_code_user_type" placeholder="Search..." />
                <datalist id="user_types"></datalist>
              </label>

              <label>ACCOUNT
                <input list="accounts" name="se_accnt" placeholder="Search..." onchange="setUserType(this, 'se_code_user_type')" />
                <datalist id="accounts"></datalist>
              </label>

              <label>ACTION BUTTON
                <input type="text" name="action_button" required placeholder="e.g., Approve" />
              </label>
              <label>ACTION FUNCTION
                <input type="text" name="action_function" placeholder="e.g., approve" />
              </label>
            </div>
            <button type="submit" class="btn primary">Create Rule</button>
          </form>
        </div>
      `;

      container.innerHTML = statusesTable + statusFormsHTML + rulesTable + ruleFormsHTML;
    })
    .catch(error => {
      console.error("Error fetching workflow details:", error);
    });


    fetch(`/lookups/${workflowID}`)
    .then(response => {
      if (!response.ok) throw new Error("Failed to fetch lookups");
      return response.json();
    })
    .then(lookups => {
        console.log("Lookups:", lookups);
        // ED Categories
        document.getElementById('ed_codes_cat').innerHTML =
        lookups.ed_status_codes_cat.map(o => `<option value="${o.ed_code_status_cat_id}_${escapeHtml(o.descr_en)}">${o.ed_code_status_cat_id}_${escapeHtml(o.descr_en)}</option>`).join('');
        ED_CODE_STATUS_CAT = lookups.ed_status_codes_cat;
      // ED
        document.getElementById('ed_codes').innerHTML =
        lookups.ed_status_codes.map(o => `<option value="${o.ed_code_status_id}_${escapeHtml(o.descr_en)}">${o.ed_code_status_id}_${escapeHtml(o.descr_en)}</option>`).join('');
        ED_CODE_STATUS = lookups.ed_status_codes;
      // GS
        document.getElementById('gs_codes').innerHTML =
        lookups.gs_status_codes.map(o => `<option value="${o.gs_code_req_status_id}_${escapeHtml(o.descr_en)}">${o.gs_code_req_status_id}_${escapeHtml(o.descr_en)}</option>`).join('');
        GS_CODE_REQ_STATUS = lookups.gs_status_codes;
        // User types
        document.getElementById('user_types').innerHTML =
        lookups.user_types.map(o => `<option value="${o.se_code_user_type_id}_${escapeHtml(o.descr_en)}">${o.se_code_user_type_id}_${escapeHtml(o.descr_en)}</option>`).join('');
        SE_CODE_USER_TYPE = lookups.user_types;
        console.log("SE_CODE_USER_TYPE:", SE_CODE_USER_TYPE);
        // Accounts
        document.getElementById('accounts').innerHTML =
        lookups.accounts.map(o => `<option value="${o.se_accnt_id}_${escapeHtml(o.descr_en)}">${o.se_accnt_id}_${escapeHtml(o.descr_en)}</option>`).join('');
        SE_ACCNT = lookups.accounts;
        console.log("SE_ACCNT:", SE_ACCNT);
        // WF statuses (from/to for rules)
        if (lookups.workflow_statuses) {
            document.getElementById('wf_statuses').innerHTML =
            lookups.workflow_statuses.map(s => `<option value="${s.status_id}_${escapeHtml(s.status_name)}">${s.status_id}_${escapeHtml(s.status_name)}</option>`).join('');
            WF_STATUSES = lookups.workflow_statuses;
        }
    })
    .catch(error => {
      console.error("Error fetching lookups:", error);
    });

}

window.setUserType = setUserType;
export function setUserType(accountInput, userTypeInputId) {
  const accountValue = accountInput.value;
  console.log("Account selected:", accountValue);
  const userTypeInput = document.getElementById(userTypeInputId);
  console.log("User Type input found:", userTypeInput);

  const accountId = accountValue.split('_')[0];
  const account = SE_ACCNT.find(acc => String(acc.se_accnt_id) === accountId);
  if (account) {
    const userType = SE_CODE_USER_TYPE.find(ut => ut.se_code_user_type_id === account.se_code_user_type_id);
    if (userType) {
      userTypeInput.value = `${userType.se_code_user_type_id}_${userType.descr_en}`;
    } else {
      userTypeInput.value = '';
    }
  } else {
    userTypeInput.value = '';
  }
}

