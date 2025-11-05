window.addEventListener('DOMContentLoaded', () => {
  const params = new URLSearchParams(window.location.search);
  const openID = params.get('workflow');
  loadLookups();
  if (openID) {
    viewDetails(openID);
  }
});

