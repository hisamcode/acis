/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./ui/htmx/**/*.html"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
}