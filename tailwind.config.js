/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: "class",
  content: ["./templates/**/*.{html,gotmpl,tmpl}",],
  theme: {
    extend: {
      colors: {
        app: {
          bg: "rgb(var(--app-bg) / <alpha-value>)",
          surface: "rgb(var(--app-surface) / <alpha-value>)",
          panel: "rgb(var(--app-panel) / <alpha-value>)",
          border: "rgb(var(--app-border) / <alpha-value>)",
          text: "rgb(var(--app-text) / <alpha-value>)",
          muted: "rgb(var(--app-muted) / <alpha-value>)",
          subtle: "rgb(var(--app-subtle) / <alpha-value>)",
          accent: "rgb(var(--app-accent) / <alpha-value>)",
          "accent-strong": "rgb(var(--app-accent-strong) / <alpha-value>)",
          danger: "rgb(var(--app-danger) / <alpha-value>)",
        },
      },
      fontFamily: {
        sans: [
          "ui-sans-serif",
          "system-ui",
          "-apple-system",
          "Segoe UI",
          "Roboto",
          "Noto Sans Arabic",
          "Tahoma",
          "Arial",
          "sans-serif",
        ],
      },
    },
  },
  plugins: [
    require("@tailwindcss/forms"),
    require("@tailwindcss/typography"),
  ],
};
