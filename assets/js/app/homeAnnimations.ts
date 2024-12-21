import "gsap";

if (typeof gsap !== "undefined") {
  gsap.registerPlugin(ScrollTrigger);
  for (let i = 1; i <= 3; i++) {
    gsap.to(`#home-ilustracion-${i}`, {
      scrollTrigger: {
        trigger: `#home-ilustracion-${i}`,
        toggleActions: "play pause resume reset",
        start: "top 80%",
      },
      duration: 1,
      opacity: 1,
    });
  }
}
