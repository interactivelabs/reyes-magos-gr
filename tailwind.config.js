/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.{html,js,templ,go}", "./asseets/css/main.css"],
  theme: {
    extend: {
      fontFamily: {
        sans: [
          "Cabin",
          "ui-sans-serif",
          "system-ui",
          "sans-serif",
          "Apple Color Emoji",
          "Segoe UI Emoji",
          "Segoe UI Symbol",
          "Noto Color Emoji",
        ],
        display: ["Raleway", "ui-sans-serif", "system-ui", "sans-serif"],
        decorative: ["Indie Flower", "cursive"],
      },
      colors: {
        "brand-orange": "#F4673D",
        "brand-blue": "#2A3245",
      },
    },
  },
  plugins: [
    require("@tailwindcss/aspect-ratio"),
    require("tailwindcss-motion"),
  ],
};
