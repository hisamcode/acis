/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./ui/htmx/**/*.html"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  safelist: [
    "toast",
    "toast-top",
    "toast-center",
    "toast-center",
    "alert",
    "alert-info"
  ]
}