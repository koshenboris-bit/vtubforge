async function loginUser() {
  const login = document.getElementById("loginInput").value.trim();
  const password = document.getElementById("passwordInput").value.trim();

  if (!login || !password) {
    return showToast(t("toastMissingData"), t("toastLoginPassword"), "warn");
  }

  const data = await apiRequest("/login", "POST", { login, password }, { allowRefresh: false });

  if (data?.error) {
    return showToast(t("toastLoginFailed"), data.error, "warn");
  }

  localStorage.setItem("accessToken", data.accessToken);
  localStorage.setItem("refreshToken", data.refreshToken);
  localStorage.setItem("user", JSON.stringify(data.user));

  showToast(t("toastWelcomeBack"), t("toastLoginSuccess"), "success");
  setTimeout(() => window.location.href = "lessons.html", 350);
}

async function registerUser() {
  const login = document.getElementById("loginInput").value.trim();
  const password = document.getElementById("passwordInput").value.trim();

  if (!login || !password) {
    return showToast(t("toastMissingData"), t("toastLoginPassword"), "warn");
  }

  const data = await apiRequest("/register", "POST", { login, password }, { allowRefresh: false });

  if (data?.error) {
    return showToast(t("toastRegistrationFailed"), data.error, "warn");
  }

  localStorage.setItem("accessToken", data.accessToken);
  localStorage.setItem("refreshToken", data.refreshToken);
  localStorage.setItem("user", JSON.stringify(data.user));

  showToast(t("toastAccountCreated"), t("toastRegisterSuccess"), "success");
  setTimeout(() => window.location.href = "lessons.html", 350);
}
