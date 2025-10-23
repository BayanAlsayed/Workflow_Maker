window.deleteWorkflow = deleteWorkflow;
export function deleteWorkflow(workflowID) {
  // Implement delete functionality here
  fetch(`/delete_workflow?Workflow_ID=${workflowID}`)
  .then(response => {
    if (response.ok) {
      console.log("Workflow deleted:", workflowID);
      // Optionally, refresh the page or remove the workflow from the UI
    } else {
      console.error("Failed to delete workflow:", workflowID);
    }
  })
  .catch(error => {
    console.error("Error deleting workflow:", error);
  });
}

window.editWorkflow = editWorkflow;
export function editWorkflow(workflowID, workflowName, workflowDescription) {
  const item_info = document.getElementById(`workflow_${workflowID}`);
  item_info.innerHTML = `
    <div class="item-info" id="workflow_info_${workflowID}">
      <input type="text" id="edit_name_${workflowID}" value="${workflowName}"/>
      <input type="text" id="edit_desc_${workflowID}" value="${workflowDescription}" />
    </div>
    <button class="icon-btn save" type="button"
      onclick="event.stopPropagation(); saveWorkflow(${workflowID})"
      title="Save">
      <i class="fa-solid fa-check"></i>
    </button>

    <button class="icon-btn cancel" type="button"
      onclick="event.stopPropagation(); cancelEdit(${workflowID}, '${workflowName}', '${workflowDescription}')"
      title="Cancel">
      <i class="fa-solid fa-xmark"></i>
    </button>

  `;
  console.log("Edit workflow:", workflowID, workflowName, workflowDescription);
}

window.saveWorkflow = saveWorkflow;
export function saveWorkflow(workflowID) {
  const newName = document.getElementById(`edit_name_${workflowID}`).value;
  const newDesc = document.getElementById(`edit_desc_${workflowID}`).value;
  // Implement save functionality here, e.g., send updated data to the server
  console.log("Save workflow:", workflowID, newName, newDesc);
  fetch(`/edit_workflow`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      workflow_id: workflowID,
      workflow_name: newName,
      workflow_description: newDesc
    }),
  })
  .then(response => {
    if (response.ok) {
      console.log("Workflow updated:", workflowID);
      // Optionally, refresh the page or update the UI
      location.reload();
    } else {
      console.error("Failed to update workflow:", workflowID);
    }
  })
  .catch(error => {
    console.error("Error updating workflow:", error);
  });
}

window.cancelEdit = cancelEdit;
export function cancelEdit(workflowID, originalName, originalDesc) {
  const item_info = document.getElementById(`workflow_${workflowID}`);
  item_info.innerHTML = `
    <div class="item-info" id="workflow_info_${workflowID}">
      <div class="item-name">${originalName}</div>
      <div class="item-desc">${originalDesc}</div>
    </div>
    <button class="icon-btn edit" type="button"
      onclick="event.stopPropagation(); editWorkflow(${workflowID}, '${originalName}', '${originalDesc}')"
      title="Edit">
      <i class="fa-solid fa-pen"></i>
    </button>

    <button class="icon-btn delete" type="button"
      onclick="event.stopPropagation(); deleteWorkflow(${workflowID})"
      title="Delete">
      <i class="fa-solid fa-trash"></i>
    </button>
  `;
  console.log("Cancel edit for workflow:", workflowID);
}
