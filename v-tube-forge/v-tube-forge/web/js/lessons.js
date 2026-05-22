function lessonTypeLabel(type) {
  return type === "full" ? t("lessonFullCourse") : t("lessonIntro");
}

function lessonTypeClass(type) {
  return type === "full" ? "full" : "intro";
}

async function passLesson(lessonId) {
  const user = getStoredUser();
  if (!user?.id) {
    return showToast(t("toastNotLoggedIn"), t("toastSignInAgain"), "warn");
  }

  const result = await apiRequest(`/lessons/${lessonId}/pass/${user.id}`, "POST", null, { allowRefresh: true });
  if (result?.error) {
    return showToast(t("toastCantMark"), result.error, "warn");
  }

  showToast(t("toastCompleted"), t("toastLessonPassed"), "success");
  await loadLessons();
  if (typeof loadProfileSummary === "function") await loadProfileSummary();
}

function renderLessonCard(lesson) {
  const passed = !!lesson.passed;
  return `
    <article class="lesson-card glass reveal">
      <div class="lesson-top">
        <span class="badge ${lessonTypeClass(lesson.lessonType)}">${lessonTypeLabel(lesson.lessonType)}</span>
        ${passed ? `<span class="badge completed">${t("lessonCompleted")}</span>` : ""}
      </div>

      <h4 class="lesson-title">${lesson.title}</h4>
      <p>${lesson.description}</p>

      <div class="lesson-actions">
        <a class="btn btn-secondary btn-sm" href="${lesson.videoLink}" target="_blank" rel="noopener">
          ${t("lessonWatch")}
        </a>
        ${
          passed
            ? `<span class="lesson-progress">${t("lessonAlreadyCompleted")}</span>`
            : `<button class="btn btn-primary btn-sm" onclick="passLesson(${lesson.id})">${t("lessonMarkCompleted")}</button>`
        }
      </div>
    </article>
  `;
}

async function loadLessons() {
  const response = await apiRequest("/lessons");

  if (response?.error) {
    showToast(t("toastError"), response.error, "warn");
    return;
  }

  const lessons = Array.isArray(response) ? response : response.data || [];
  const introEl = document.getElementById("introLessons");
  const fullEl = document.getElementById("fullLessons");

  if (!introEl || !fullEl) return;

  const introLessons = lessons.filter(item => item.lessonType === "intro");
  const fullLessons = lessons.filter(item => item.lessonType === "full");

  introEl.innerHTML = introLessons.length
    ? introLessons.map(renderLessonCard).join("")
    : `<div class="notice">${t("noIntroLessons")}</div>`;

  fullEl.innerHTML = fullLessons.length
    ? fullLessons.map(renderLessonCard).join("")
    : `<div class="notice">${t("noFullLessons")}</div>`;
}

async function loadLessonsAdmin() {
  const data = await apiRequest("/lessons");
  if (data?.error) return [];

  const list = document.getElementById("adminLessonsList");
  if (!list) return data || [];

  list.innerHTML = (data || []).map(item => `
    <div class="panel-item">
      <div>
        <strong>${item.title}</strong>
        <small>${lessonTypeLabel(item.lessonType)} · ${item.description}</small>
      </div>
      <button class="btn btn-danger btn-sm" onclick="deleteLesson(${item.id})">${t("adminDelete")}</button>
    </div>
  `).join("") || `<div class="notice">${t("adminNoLessons")}</div>`;

  return data || [];
}

async function createLesson() {
  const title = document.getElementById("lessonTitle").value.trim();
  const description = document.getElementById("lessonDescription").value.trim();
  const lessonType = document.getElementById("lessonType").value;
  const videoLink = document.getElementById("lessonVideo").value.trim();

  if (!title || !description || !videoLink) {
    return showToast(t("toastMissingData"), t("toastFillLesson"), "warn");
  }

  const result = await apiRequest("/lessons", "POST", { title, description, lessonType, videoLink });
  if (result?.error) return showToast(t("toastError"), result.error, "warn");

  showToast(t("toastLessonCreated"), t("toastLessonAdded"), "success");
  document.getElementById("lessonTitle").value = "";
  document.getElementById("lessonDescription").value = "";
  document.getElementById("lessonVideo").value = "";

  await loadLessonsAdmin();
}

async function deleteLesson(id) {
  if (!confirm(t("confirmDeleteLesson"))) return;
  const result = await apiRequest(`/lessons/${id}`, "DELETE");
  if (result?.error) return showToast(t("toastError"), result.error, "warn");
  showToast(t("toastLessonDeleted"), t("toastRemoved"), "success");
  await loadLessonsAdmin();
  if (document.getElementById("introLessons")) await loadLessons();
}

window.addEventListener("vtuberforge:languagechange", () => {
  if (document.getElementById("introLessons")) loadLessons();
  if (document.getElementById("adminLessonsList")) loadLessonsAdmin();
});
