document.addEventListener("alpine:init", () => {
  Alpine.data("app", () => ({
    theme: "light",
    lang: "en",

    init() {
      this.theme = document.documentElement.classList.contains("dark") ? "dark" : "light";
      this.lang = document.documentElement.lang || "en";
    },

    applyTheme(theme) {
      this.theme = theme;
      document.documentElement.classList.toggle("dark", theme === "dark");
      try {
        localStorage.setItem("theme", theme);
      } catch {}
    },

    toggleTheme() {
      this.applyTheme(this.theme === "dark" ? "light" : "dark");
    },

    async setLang(lang) {
      if (lang === this.lang) {
        return;
      }

      const previousLang = this.lang;
      const previousDir = document.documentElement.dir;
      document.documentElement.lang = lang;
      document.documentElement.dir = lang === "ar" ? "rtl" : "ltr";

      try {
        const response = await fetch("/language", {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
          body: new URLSearchParams({ lang }),
        });

        if (response.ok) {
          window.location.reload();
          return;
        }
      } catch {}

      document.documentElement.lang = previousLang;
      document.documentElement.dir = previousDir;
    },
  }));
});
