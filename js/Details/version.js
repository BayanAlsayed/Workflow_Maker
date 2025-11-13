window.activateWorkflowVersion = activateWorkflowVersion;
export function activateWorkflowVersion(workflowID, version) {
    const data = {
        workflow_id: workflowID,
        version: version
    };

    fetch(`/activate_workflow_version`, 
        { method: 'POST', 
        headers: { 'Content-Type': 'application/json' }, 
        body: JSON.stringify(data) })
    .then(r => { if (!r.ok) throw new Error('Failed to activate workflow version'); })
    .then(() => {
        viewDetails(workflowID, version);
    });
}

window.approveWorkflowVersion = approveWorkflowVersion;
export function approveWorkflowVersion(workflowID, version) {
    const data = {
        workflow_id: workflowID,
        version: version
    };

    fetch(`/approve_workflow_version`, 
        { method: 'POST', 
        headers: { 'Content-Type': 'application/json' }, 
        body: JSON.stringify(data) })
    .then(r => { if (!r.ok) throw new Error('Failed to activate workflow version'); })
    .then(() => {
        viewDetails(workflowID, version);
    });
}

window.deleteWorkflowVersion = deleteWorkflowVersion;
export function deleteWorkflowVersion(workflowID, version) {
    const data = {
        workflow_id: workflowID,
        version: version
    };

    fetch(`/delete_workflow_version`, 
        { method: 'POST', 
        headers: { 'Content-Type': 'application/json' }, 
        body: JSON.stringify(data) })
    .then(r => { if (!r.ok) throw new Error('Failed to activate workflow version'); })
    .then(() => {
        viewDetails(workflowID, 0);
    });
}

window.createWorkflowVersion = createWorkflowVersion;
export function createWorkflowVersion(workflowID) {
    fetch(`/create_workflow_version`, 
        { method: 'POST', 
        headers: { 'Content-Type': 'application/json' }, 
        body: JSON.stringify(workflowID) })
    .then(r => { if (!r.ok) throw new Error('Failed to activate workflow version'); })
    .then(version => {
        viewDetails(workflowID, version);
    });
}

window.duplicateWorkflowVersion = duplicateWorkflowVersion;
export function duplicateWorkflowVersion(workflowID, version) {
    const data = {
        workflow_id: workflowID,
        version: version
    }
    fetch(`/duplicate_workflow_version`, 
        { method: 'POST', 
        headers: { 'Content-Type': 'application/json' }, 
        body: JSON.stringify(data) })
    .then(r => { if (!r.ok) throw new Error('Failed to activate workflow version'); })
    .then(version => {
        viewDetails(workflowID, version);
    });
}

