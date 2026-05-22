function getStoredUser() {
  const raw = localStorage.getItem("user");
  if (!raw) return null;
  try { return JSON.parse(raw); } catch { return null; }
}

function isLoggedIn() {
  return !!localStorage.getItem("accessToken");
}

function isAdmin() {
  const user = getStoredUser();
  return !!user && user.role === "admin";
}

function logout() {
  localStorage.removeItem("accessToken");
  localStorage.removeItem("refreshToken");
  localStorage.removeItem("user");
  window.location.href = "index.html";
}

function publicPage() {
  const page = document.body?.dataset.page || "";
  return ["home", "login", "register", "try-vtuber"].includes(page);
}

function shouldProtectAdmin() {
  return document.body?.dataset.page === "admin";
}

function shouldProtectPrivate() {
  return ["lessons", "news", "profile", "admin"].includes(document.body?.dataset.page || "");
}

function showToast(title, message, type = "info") {
  const existing = document.querySelector(".toast");
  if (existing) existing.remove();

  const el = document.createElement("div");
  el.className = `toast ${type}`;
  el.innerHTML = `
    <strong>${title}</strong>
    <div class="muted">${message || ""}</div>
  `;
  document.body.appendChild(el);
  requestAnimationFrame(() => el.classList.add("show"));
  setTimeout(() => {
    el.classList.remove("show");
    setTimeout(() => el.remove(), 250);
  }, 2600);
}

function apiHeaders(includeAuth = true) {
  const headers = { "Content-Type": "application/json" };
  if (includeAuth) {
    const token = localStorage.getItem("accessToken");
    if (token) headers["Authorization"] = `Bearer ${token}`;
  }
  return headers;
}

async function parseResponse(response) {
  if (response.status === 204) return null;
  const contentType = response.headers.get("content-type") || "";
  if (contentType.includes("application/json")) return await response.json();
  return { raw: await response.text() };
}

async function fetchOnBase(baseUrl, endpoint, method, body, includeAuth = true) {
  const url = `${baseUrl}${endpoint}`;
  return fetch(url, {
    method,
    headers: apiHeaders(includeAuth),
    body: body ? JSON.stringify(body) : null
  });
}

async function refreshToken() {
  const refreshToken = localStorage.getItem("refreshToken");
  if (!refreshToken) return false;

  for (const baseUrl of (window.APP_CONFIG?.apiBases || [])) {
    try {
      const response = await fetchOnBase(baseUrl, "/refresh", "POST", { refreshToken }, false);
      if (response.status === 404) continue;
      if (!response.ok) continue;

      const data = await parseResponse(response);
      if (!data || !data.accessToken || !data.refreshToken) continue;

      localStorage.setItem("accessToken", data.accessToken);
      localStorage.setItem("refreshToken", data.refreshToken);
      return true;
    } catch (_) {}
  }

  logout();
  return false;
}

async function apiRequest(endpoint, method = "GET", body = null, options = {}) {
  const bases = window.APP_CONFIG?.apiBases;
  const includeAuth = options.includeAuth !== false;
  const allowRefresh = options.allowRefresh !== false;

  let lastError = null;

  for (let i = 0; i < bases.length; i++) {
    const baseUrl = bases[i];
    try {
      const response = await fetchOnBase(baseUrl, endpoint, method, body, includeAuth);

      if (response.status === 404 && i < bases.length - 1) {
        continue;
      }

      if (response.status === 401 && allowRefresh) {
        const refreshed = await refreshToken();
        if (refreshed) {
          return apiRequest(endpoint, method, body, { ...options, allowRefresh: false });
        }
      }

      const data = await parseResponse(response);

      if (!response.ok) {
        if (data && typeof data === "object" && data.error) return data;
        return { error: `Request failed with status ${response.status}` };
      }

      return data;
    } catch (err) {
      lastError = err;
      continue;
    }
  }

  return {
    error: lastError?.message || "Unable to reach backend"
  };
}

function setNavState() {
  const navAuth = document.querySelectorAll("[data-nav='auth']");
  const navUser = document.querySelectorAll("[data-nav='user']");
  const navAdmin = document.querySelectorAll("[data-nav='admin']");
  const navLogout = document.querySelectorAll("[data-nav='logout']");
  const navLogin = document.querySelectorAll("[data-nav='login']");
  const navRegister = document.querySelectorAll("[data-nav='register']");

  const loggedIn = isLoggedIn();
  const admin = isAdmin();

  navAuth.forEach(el => el.classList.toggle("hidden", loggedIn));
  navUser.forEach(el => el.classList.toggle("hidden", !loggedIn));
  navLogout.forEach(el => el.classList.toggle("hidden", !loggedIn));
  navAdmin.forEach(el => el.classList.toggle("hidden", !(loggedIn && admin)));
  navLogin.forEach(el => el.classList.toggle("hidden", loggedIn));
  navRegister.forEach(el => el.classList.toggle("hidden", loggedIn));
}

function protectRoute() {
  const page = document.body?.dataset.page || "";
  const loggedIn = isLoggedIn();
  const admin = isAdmin();

  if (publicPage()) return;

  if (!loggedIn && shouldProtectPrivate()) {
    window.location.href = "login.html";
    return;
  }

  if (page === "admin" && !admin) {
    window.location.href = "index.html";
    return;
  }

  if ((page === "login" || page === "register") && loggedIn) {
    window.location.href = "lessons.html";
  }
}

function setupScrollState() {
  const nav = document.querySelector(".nav");
  if (!nav) return;

  const update = () => {
    nav.classList.toggle("scrolled", window.scrollY > 12);
  };

  update();
  window.addEventListener("scroll", update, { passive: true });
}

function setupReveal() {
  const els = document.querySelectorAll(".reveal");
  if (!els.length) return;

  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        entry.target.classList.add("visible");
        observer.unobserve(entry.target);
      }
    });
  }, { threshold: 0.12 });

  els.forEach(el => observer.observe(el));
}

function setupParticles() {
  const canvas = document.getElementById("particles-canvas");
  if (!canvas) return;

  const ctx = canvas.getContext("2d");
  let width = 0;
  let height = 0;
  let particles = [];
  const count = 58;

  const colors = [
    ["124,58,237", 0.45],
    ["56,189,248", 0.35],
    ["244,114,182", 0.28]
  ];

  const resize = () => {
    width = canvas.width = window.innerWidth;
    height = canvas.height = window.innerHeight;
  };

  const create = () => {
    const [color, alpha] = colors[Math.floor(Math.random() * colors.length)];
    return {
      x: Math.random() * width,
      y: Math.random() * height,
      r: 0.8 + Math.random() * 1.6,
      vx: (Math.random() - 0.5) * 0.35,
      vy: (Math.random() - 0.5) * 0.35,
      color,
      alpha
    };
  };

  const init = () => {
    resize();
    particles = Array.from({ length: count }, create);
  };

  const draw = () => {
    ctx.clearRect(0, 0, width, height);

    for (const p of particles) {
      p.x += p.vx;
      p.y += p.vy;

      if (p.x < 0) p.x = width;
      if (p.x > width) p.x = 0;
      if (p.y < 0) p.y = height;
      if (p.y > height) p.y = 0;

      ctx.beginPath();
      ctx.arc(p.x, p.y, p.r, 0, Math.PI * 2);
      ctx.fillStyle = `rgba(${p.color},${p.alpha})`;
      ctx.fill();
    }

    for (let i = 0; i < particles.length; i++) {
      for (let j = i + 1; j < particles.length; j++) {
        const a = particles[i];
        const b = particles[j];
        const dist = Math.hypot(a.x - b.x, a.y - b.y);
        if (dist < 120) {
          ctx.beginPath();
          ctx.moveTo(a.x, a.y);
          ctx.lineTo(b.x, b.y);
          ctx.strokeStyle = `rgba(124,58,237,${0.08 * (1 - dist / 120)})`;
          ctx.lineWidth = 0.6;
          ctx.stroke();
        }
      }
    }

    requestAnimationFrame(draw);
  };

  window.addEventListener("resize", resize);
  init();
  draw();
}

document.addEventListener("DOMContentLoaded", () => {
  setNavState();
  protectRoute();
  setupScrollState();
  setupReveal();
  setupParticles();

  document.querySelectorAll("[data-logout]").forEach(btn => {
    btn.addEventListener("click", logout);
  });

  document.querySelectorAll("[data-nav-link]").forEach(link => {
    link.addEventListener("click", () => {
      document.querySelectorAll("[data-nav-link]").forEach(a => a.classList.remove("active"));
      link.classList.add("active");
    });
  });
});

window.VTuberApp = {
  logout,
  apiRequest,
  isLoggedIn,
  isAdmin,
  getStoredUser,
  showToast
};

