async function initAdmin() {
  if (!isAdmin()) {
    return;
  }
  await loadLessonsAdmin();
  await loadNews();
}
