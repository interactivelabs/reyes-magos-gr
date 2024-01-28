/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.{html,js,templ,go}"],
  theme: {
    extend: {
      backgroundImage: {
        "hero-paper": "url('/public/img/bg_water02_sm.webp')",
        "footer-texture": "url('/public/img/footer_bg.webp')",
      },
      fontFamily: {
        display: ['"Indie Flower"', "cursive"],
      },
      colors: {
        "brand-orange": "#F4673D",
        "brand-blue": "#2A3245",
      },
    },
  },
  plugins: [],
};
