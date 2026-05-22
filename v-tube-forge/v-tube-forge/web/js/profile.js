async function loadProfileSummary() {
  const me = await apiRequest("/users/me");
  if (me?.error) {
    showToast(t("toastProfileError"), me.error, "warn");
    return;
  }

  const profileName = document.getElementById("profileName");
  const profileRole = document.getElementById("profileRole");
  const lastLesson = document.getElementById("lastLesson");
  const completedCount = document.getElementById("completedCount");
  const completedList = document.getElementById("completedList");

  if (!profileName || !profileRole || !lastLesson || !completedCount || !completedList) return;

  profileName.textContent = me.login;
  profileRole.textContent = me.role;
  lastLesson.textContent = me.lastLesson ? me.lastLesson.title : t("profileNoCompleted");

  const lessons = await apiRequest("/lessons");
  const completed = (lessons || []).filter(item => item.passed);

  completedCount.textContent = String(completed.length);
  completedList.innerHTML = completed.length
    ? completed.map(item => `<span class="chip">${item.title}</span>`).join("")
    : `<div class="notice">${t("profileNoLessons")}</div>`;
}

window.addEventListener("vtuberforge:languagechange", () => {
  if (document.getElementById("profileName")) loadProfileSummary();
});
