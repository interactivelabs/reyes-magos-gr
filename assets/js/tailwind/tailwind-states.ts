import { addStateClasses, removeStateClasses } from "./tailwind-core";

export default function initTailwindStates() {
  const ATTRIBUTE = "data-styles-state";

  const mutationCallback = (mutationsList) => {
    for (const mutation of mutationsList) {
      const { type, attributeName, oldValue, element } = mutation;
      if (type !== "attributes" || attributeName !== ATTRIBUTE || !element) {
        return;
      }

      // Remove previous state classes
      if (oldValue) {
        removeStateClasses(element, `data-styles-${oldValue}`);
      }

      const state = element.getAttribute(ATTRIBUTE);
      addStateClasses(element, `data-styles-${state}`);
    }
  };

  const observer = new MutationObserver(mutationCallback);

  document.addEventListener("DOMContentLoaded", function () {
    const allElements = document.querySelectorAll(`[${ATTRIBUTE}]`);

    allElements.forEach((transitionElement: Element) => {
      const state = transitionElement.getAttribute(ATTRIBUTE);
      addStateClasses(transitionElement, `data-styles-${state}`);
      observer.observe(transitionElement, {
        attributes: true,
        attributeOldValue: true,
        attributeFilter: [ATTRIBUTE],
      });
    });
  });
}
