async function loadNews() {
  const data = await apiRequest("/news");
  const news = Array.isArray(data) ? data : [];

  const featured = document.getElementById("featuredNews");
  const list = document.getElementById("newsList");
  const adminList = document.getElementById("adminNewsList");

  if (featured) {
    const top = news[0];
    featured.innerHTML = top ? `
      <div class="news-hero-main">
        <span class="news-eyebrow">${t("featuredUpdate")}</span>
        <h2 class="news-hero-title">${top.title}</h2>
        <p class="news-hero-text">${top.content}</p>
      </div>
      <div class="news-hero-side">
        <div class="news-feature-list">
          <div class="news-feature">${t("newsFeature1")}</div>
          <div class="news-feature">${t("newsFeature2")}</div>
          <div class="news-feature">${t("newsFeature3")}</div>
        </div>
      </div>
    ` : `<div class="notice">${t("noNews")}</div>`;
  }

  if (list) {
    const items = news.slice(1);
    list.innerHTML = items.length
      ? items.map((item, index) => `
        <article class="news-card glass reveal">
          <span class="news-date">${t("newsUpdate")} ${index + 2}</span>
          <h3>${item.title}</h3>
          <p>${item.content}</p>
          <div class="news-footer">
            <span>${t("newsBulletin")}</span>
            <span>${t("newsLive")}</span>
          </div>
        </article>
      `).join("")
      : `<div class="notice">${t("noMoreNews")}</div>`;
  }

  if (adminList) {
    adminList.innerHTML = news.length
      ? news.map(item => `
        <div class="panel-item">
          <div>
            <strong>${item.title}</strong>
            <small>${item.content}</small>
          </div>
          <button class="btn btn-danger btn-sm" onclick="deleteNews(${item.id})">${t("adminDelete")}</button>
        </div>
      `).join("")
      : `<div class="notice">${t("adminNoNews")}</div>`;
  }
}

async function createNews() {
  const title = document.getElementById("newsTitle").value.trim();
  const content = document.getElementById("newsContent").value.trim();

  if (!title || !content) {
    return showToast(t("toastMissingData"), t("toastFillNews"), "warn");
  }

  const result = await apiRequest("/news", "POST", { title, content });
  if (result?.error) return showToast(t("toastError"), result.error, "warn");

  showToast(t("toastNewsCreated"), t("toastNewsAdded"), "success");
  document.getElementById("newsTitle").value = "";
  document.getElementById("newsContent").value = "";
  await loadNews();
}

async function deleteNews(id) {
  if (!confirm(t("confirmDeleteNews"))) return;
  const result = await apiRequest(`/news/${id}`, "DELETE");
  if (result?.error) return showToast(t("toastError"), result.error, "warn");
  showToast(t("toastNewsDeleted"), t("toastRemoved"), "success");
  await loadNews();
}

window.addEventListener("vtuberforge:languagechange", () => {
  if (document.getElementById("featuredNews") || document.getElementById("newsList") || document.getElementById("adminNewsList")) loadNews();
});
