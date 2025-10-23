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


