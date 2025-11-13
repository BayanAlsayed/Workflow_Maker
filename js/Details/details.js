// details.js
import { escapeHtml, nullable, opt, buildOptions, initTomSelectSingle } from "../helpers.js";
export let ED_CODE_STATUS_CAT = [];
export let ED_CODE_STATUS = [];
export let GS_CODE_REQ_STATUS = [];
export let SE_CODE_USER_TYPE = [];
export let SE_ACCNT = [];
export let WF_CONDITIONS = [];
export let WF_STATUSES = [];



window.viewDetails = viewDetails;
export function viewDetails(workflowID, version) {
  

  const container = document.getElementById('main_container');
  container.innerHTML = '';
  container.style.display = 'block';

  fetch(`/get_workflow_versions/${workflowID}`)
    .then(r => { if (!r.ok) throw new Error("Failed to fetch workflow version"); return r.json(); })
    .then(versions => {
      if (version > 0) {
        history.pushState({}, '', `?workflow=${workflowID}&version=${version}`);     
      } else if (versions.filter(v => v.is_active).length > 0){
        version = versions.find(v => v.is_active).version;
        history.pushState({}, '', `?workflow=${workflowID}&version=${versions.find(v => v.is_active).version}`);     
      } else {
        version = versions[0].version
        history.pushState({}, '', `?workflow=${workflowID}&version=${versions[0].version}`);     
      }

      if (versions.length > 0) {
        const versionsSelect = `
          <label for="version_sel">Version:
            <select id="version_sel" onchange="viewDetails(${workflowID}, this.value)">
              ${versions.map(v => 
                `<option value="${v.version}" ${v.version == version ? 'selected' : ''}>
                  ${v.version} ${v.is_active ? '(Active)' : v.is_approved ? '(Approved)' : ''}
                </option>`).join('')}
            </select>
          </label>
          <button id="create_workflow_version_button" onclick="createWorkflowVersion(${workflowID})">create version</button>
          <button id="duplicate_workflow_version_button" onclick="duplicateWorkflowVersion(${workflowID}, ${version})">duplicate version</button>
        `;

        container.innerHTML = versionsSelect;
        // var version = document.getElementById("version_sel").value || versions[0].version

        const versionData = versions?.find(v => String(v.version) === String(version));
        if (versionData && versionData.is_active) {
          // active version - no activate button
        } else  if (versionData && versionData.is_approved) {
          const versionActivateBtn = `
            <button id="activate_version_btn" class="btn primary" onclick="activateWorkflowVersion(${workflowID}, ${version})">
              Activate Version
            </button>
          `;
          container.insertAdjacentHTML('beforeend', versionActivateBtn);
        } else  if (versionData) {
          const versionApproveBtn = `
            <button id="approve_version_btn" class="btn primary" onclick="approveWorkflowVersion(${workflowID}, ${version})">
              Approve Version
            </button>
            <button id="delete_version_btn" class="btn primary" onclick="deleteWorkflowVersion(${workflowID}, ${version})">
              Delete Version
            </button>
          `;
          container.insertAdjacentHTML('beforeend', versionApproveBtn);
        } else {
          console.warn('Version not found in versions list:', version);
        }

        fetch(`/view_workflow/${workflowID}/${version}`)
          .then(r => { if (!r.ok) throw new Error("Failed to fetch workflow details"); return r.json(); })
          .then(data => {
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
                    <tr><td colspan="7">No statuses found.</td></tr>
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
                          onclick="event.stopPropagation(); editStatus(${s.status_id}, '${escapeHtml(s.status_name)}', ${s.ed_code_status_id ?? 'null'}, '${escapeHtml(nullable(s.ed_descr_en)||'')}', ${s.gs_code_req_status_id ?? 'null'}, '${escapeHtml(nullable(s.gs_descr_en)||'')}', ${s.is_terminal?1:0}, ${s.success_path==null?'null':Number(s.success_path)}, ${workflowID}, ${version})"
                          title="Edit"
                          ${versionData && versionData.is_approved ? 'disabled' : ''}>
                          <i class="fa-solid fa-pen"></i>
                        </button>

                        <button class="icon-btn delete" type="button"
                          onclick="event.stopPropagation(); deleteStatus(${s.status_id}, ${workflowID}, ${version})"
                          title="Delete"
                          ${versionData && versionData.is_approved ? 'disabled' : ''}>
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
                          onclick="event.stopPropagation(); editRule(${r.rule_id}, ${r.from_status_id}, '${escapeHtml(nullable(r.from_status_name)||'')}', ${r.to_status_id}, '${escapeHtml(nullable(r.to_status_name)||'')}', ${r.se_code_user_type_id ?? 'null'}, '${escapeHtml(nullable(r.user_type_en)||'')}', ${r.se_accnt_id ?? 'null'}, '${escapeHtml(nullable(r.accnt_en)||'')}', '${escapeHtml(r.action_button||'')}', '${escapeHtml(r.action_function||'')}', ${workflowID}, ${version}, ${r.rule_conditions && r.rule_conditions.length ? true : false})"
                          title="Edit"
                          ${versionData && versionData.is_approved ? 'disabled' : ''}>
                          <i class="fa-solid fa-pen"></i>
                        </button>

                        <button class="icon-btn delete" type="button"
                          onclick="event.stopPropagation(); deleteRule(${r.rule_id}, ${workflowID}, ${version})"
                          title="Delete"
                          ${versionData && versionData.is_approved ? 'disabled' : ''}>
                          <i class="fa-solid fa-trash"></i>
                        </button>

                        <button class="icon-btn" onclick="event.stopPropagation(); toggleRuleConditions(${r.rule_id})" title="Toggle conditions">
                          <i class="fa-solid fa-plus" id="icon_plus_${r.rule_id}" style="transform: rotate(${r.rule_conditions && r.rule_conditions.length ? '45deg' : '0deg'}); transition: transform 0.2s;"></i>
                        </button>
                      </td>
                    </tr>
                    <tr id="rule_conditions_wrapper_${r.rule_id}" class="cond-row" style="display: ${r.rule_conditions && r.rule_conditions.length ? 'table-row' : 'none'};">
                      <td class="cond-cell" colspan="8">
                        <div class="cond-container" id="rule_conditions_container_${r.rule_id}">
                          ${
                            (r.rule_conditions && r.rule_conditions.length) ? 
                            r.rule_conditions.map(rc => `
                              <div class="cond-card" id="cond_card_${rc.rule_condition_id}_${r.rule_id}">
                                <div class="cond-func">${escapeHtml(rc.wf_condition_func_name || '(Unnamed)')}</div>
                                <div class="cond-desc">${escapeHtml(rc.wf_condition_description || '')}</div>
                                <div class="cond-meta">
                                  <div class="cond-left">
                                    <span>${escapeHtml(rc.rule_condition_type || '')}</span>
                                  </div>
                                  <div class="cond-actions">
                                    <button class="icon-btn cond" title="Edit condition" 
                                    onclick="event.stopPropagation(); editRuleCondition(${rc.wf_rule_condition_id}, ${rc.rule_condition_id}, ${r.rule_id}, ${rc.wf_condition_id}, '${rc.wf_condition_func_name}', '${rc.wf_condition_description}', '${rc.rule_condition_type}', ${workflowID}, ${version})"
                                    ${versionData && versionData.is_approved ? 'disabled' : ''}>
                                    <i class="fa-solid fa-pen"></i></button>
                                    <button class="icon-btn cond" title="Delete condition" 
                                    onclick="event.stopPropagation(); deleteRuleCondition(${rc.wf_rule_condition_id}, ${workflowID}, ${version})"
                                    ${versionData && versionData.is_approved ? 'disabled' : ''}>
                                    <i class="fa-solid fa-trash"></i></button>
                                  </div>
                                </div>
                              </div>
                            `).join('') :
                            `<div class="cond-empty">No conditions attached to this rule.</div>`
                          }

                          <!-- add condition quick button -->
                          <div style="align-self:flex-start; margin-left:6px;">
                            <button class="cond-add-btn" 
                            onclick="event.stopPropagation(); addRuleConditionPrompt(${r.rule_id}, ${workflowID}, ${version})"
                            ${versionData && versionData.is_approved ? 'disabled' : ''}>
                              <i class="fa-solid fa-plus-circle" style="margin-right:6px"></i>
                              Add condition
                            </button>
                          </div>

                        </div>
                      </td>
                    </tr>
                  `).join('')}
                </tbody>
              </table>
            `;

            // Add forms with <select> (Tom Select will attach)
            const statusFormsHTML = versionData && !versionData.is_approved ? `
              <div class="wf-create">
                <h3>Add Status</h3>
                <form id="add-status-form" onsubmit="createStatus(event, ${workflowID}, ${version})">
                  <div class="grid">
                    <label>STATUS_NAME
                      <input type="text" name="status_name" required />
                    </label>

                    <label>ED_CODE_STATUS_CAT
                      <select id="ed_cat_sel" name="ed_code_status_cat_id">
                        <option value="">--</option>
                      </select>
                    </label>

                    <label>ED_CODE_STATUS
                      <select id="ed_code_sel" name="ed_code_status_id">
                        <option value="">--</option>
                      </select>
                    </label>

                    <label>GS_REQ_STATUS
                      <select id="gs_sel" name="gs_code_req_status_id">
                        <option value="">--</option>
                      </select>
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
              </div>`: ``;

            const ruleFormsHTML = versionData && !versionData.is_approved ? `
              <div class="wf-create">
                <h3>Add Rule</h3>
                <form id="add-rule-form" onsubmit="createRule(event, ${workflowID}, ${version})">
                  <div class="grid">
                    <label>FROM STATUS
                      <select id="from_status_sel" name="from_status" required>
                        <option value="">--</option>
                      </select>
                    </label>

                    <label>TO STATUS
                      <select id="to_status_sel" name="to_status" required>
                        <option value="">--</option>
                      </select>
                    </label>

                    <label>USER TYPE
                      <select id="user_type_sel" name="se_code_user_type">
                        <option value="">--</option>
                      </select>
                    </label>

                    <label>ACCOUNT
                      <select id="account_sel" name="se_accnt">
                        <option value="">--</option>
                      </select>
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
              </div>`: ``;

            container.innerHTML += statusesTable + statusFormsHTML + rulesTable + ruleFormsHTML;
          })
          .catch(err => console.error("Error fetching workflow details:", err));

        // Lookups & init Tom Selects
        fetch(`/wf_lookups/${workflowID}/${version}`)
          .then(r => { if (!r.ok) throw new Error("Failed to fetch lookups"); return r.json(); })
          .then(statuses => {

            WF_STATUSES = statuses || [];

            console.log("ED_CODE_STATUS: ", ED_CODE_STATUS);

            // Populate "Add Status" selects
            const edCatSel = document.getElementById('ed_cat_sel');
            const edCodeSel = document.getElementById('ed_code_sel');
            const gsSel = document.getElementById('gs_sel');

            // populate raw selects once (before TS init)
            if (edCatSel) edCatSel.innerHTML =
              `<option value="">--</option>` +
              buildOptions(ED_CODE_STATUS_CAT, 'ed_code_status_cat_id', 'descr_en',
                o => `${o.descr_en} (ID:${o.ed_code_status_cat_id})`);

            if (edCodeSel) edCodeSel.innerHTML = `<option value="">--</option>`; // start empty

            if (gsSel) gsSel.innerHTML =
              `<option value="">--</option>` +
              buildOptions(GS_CODE_REQ_STATUS, 'gs_code_req_status_id', 'descr_en',
                o => `${o.descr_en} (ID:${o.gs_code_req_status_id})`);

            const EMPTY_OPTION = { value: '', text: '--' };
            // init Tom Select ONCE and keep references
            const tsCat = initTomSelectSingle('#ed_cat_sel');
            const tsEd  = new TomSelect('#ed_code_sel', {
              create: false,
              allowEmptyOption: true,
              closeAfterSelect: true,
              valueField: 'value',
              labelField: 'text',
              searchField: ['text'],
              options: [EMPTY_OPTION] // we’ll load after cat change
            });
            initTomSelectSingle('#gs_sel');

            // helper: get the right category key from your payload
            const catKey = (o) => String(
              o.cat_id ??
              o.ED_CODE_STATUS_CAT_ID ??
              o.ed_code_status_cat_id
            );

            // when category changes, rebuild ED codes using TS API (NO innerHTML!)
            tsCat.on('change', (val) => {
              const list = (val ? ED_CODE_STATUS.filter(x => catKey(x) === String(val)) : [])
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
            });

            // Populate "Add Rule" selects
            const fromSel = document.getElementById('from_status_sel');
            const toSel   = document.getElementById('to_status_sel');
            const utSel   = document.getElementById('user_type_sel');
            const accSel  = document.getElementById('account_sel');

            if (fromSel) fromSel.innerHTML = `<option value="">--</option>` +
              buildOptions(WF_STATUSES, 'status_id', 'status_name', o => `${o.status_name} (ID:${o.status_id})`);
            if (toSel) toSel.innerHTML = `<option value="">--</option>` +
              buildOptions(WF_STATUSES, 'status_id', 'status_name', o => `${o.status_name} (ID:${o.status_id})`);
            if (utSel) utSel.innerHTML = `<option value="">--</option>` +
              buildOptions(SE_CODE_USER_TYPE, 'se_code_user_type_id', 'descr_en', o => `${o.descr_en} (ID:${o.se_code_user_type_id})`);
            if (accSel) accSel.innerHTML = `<option value="">--</option>` +
              buildOptions(SE_ACCNT, 'se_accnt_id', 'descr_en', o => `${o.descr_en} (ID:${o.se_accnt_id})`);

            const tsFrom = initTomSelectSingle('#from_status_sel');
            const tsTo   = initTomSelectSingle('#to_status_sel');
            const tsUT   = initTomSelectSingle('#user_type_sel');
            const tsAcc  = initTomSelectSingle('#account_sel', (val) => {
              // when account changes, auto-pick its user type
              const acct = SE_ACCNT.find(a => String(a.se_accnt_id) === String(val));
              if (!acct) return;
              const ut = SE_CODE_USER_TYPE.find(u => String(u.se_code_user_type_id) === String(acct.se_code_user_type_id));
              if (ut && tsUT) tsUT.setValue(String(ut.se_code_user_type_id), true);
            });

            // expose for rule.js to import arrays
            // (already exported at top)
          })
          .catch(err => console.error("Error fetching lookups:", err));
      }
    })
    .catch(err => console.error("Error fetching workflow version:", err));
}

window.loadLookups = loadLookups;
export function loadLookups() {
  // Preload CONDITIONS for condition.js
  fetch('/lookups')
    .then(r => { if (!r.ok) throw new Error("Failed to fetch conditions lookup"); return r.json(); })
    .then(lookups => {
        ED_CODE_STATUS_CAT = lookups.ed_status_codes_cat || [];
        ED_CODE_STATUS     = lookups.ed_status_codes || [];
        GS_CODE_REQ_STATUS = lookups.gs_status_codes || [];
        SE_CODE_USER_TYPE  = lookups.user_types || [];
        SE_ACCNT           = lookups.accounts || [];
        WF_CONDITIONS      = lookups.wf_conditions || [];

        console.log("Loaded lookups:", {
          ED_CODE_STATUS_CAT,
          ED_CODE_STATUS,
          GS_CODE_REQ_STATUS,
          SE_CODE_USER_TYPE,
          SE_ACCNT,
          WF_CONDITIONS
        });
    })
    .catch(err => console.error("Error fetching conditions lookup:", err));
}


