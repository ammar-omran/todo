/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.{html,gotmpl,tmpl}",],
  theme: {
    extend: {
      colors: {
        primary: "#6366f1",
      },
    },
  },
  plugins: [
    require("@tailwindcss/forms"),
    require("@tailwindcss/typography"),
  ],
};
