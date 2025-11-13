window.addEventListener('DOMContentLoaded', () => {
  const params = new URLSearchParams(window.location.search);
  const openWFID = params.get('workflow');
  const openWFVersion = params.get('version')
  loadLookups();
  if (openWFID) {
    viewDetails(openWFID, openWFVersion);
  }
});

