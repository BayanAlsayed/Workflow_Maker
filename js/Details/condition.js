window.viewConditions = viewConditions;
export function viewConditions() {
  const container = document.getElementById('main_container');
  container.innerHTML = ''; // Clear previous content

  fetch('/view_conditions', { method: 'GET', headers: { 'Cache-Control': 'no-cache' } })
    .then(response => response.json())
    .then(data => {
      data.forEach(cond => {
        const div = document.createElement('div');
        div.className = 'condition';
        div.innerHTML = `<h4>${cond.name}</h4><p>${cond.description}</p>`;
        container.appendChild(div);
      });
    });

  container.style.display = 'block'; // Show the container
}
