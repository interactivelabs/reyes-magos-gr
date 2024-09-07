import {
  addStateClasses,
  getDelayFromAttribute,
  removeStateClasses,
} from "./tailwind-core";

export default function initTailwindTransitions() {
  const ATTRIBUTE = "data-transition-state";

  const hideElement = (element: HTMLElement) => {
    element.style.display = "none";
  };

  const showElement = (element: HTMLElement) => {
    element.style.display = "block";
    element.offsetHeight; // Trigger reflow
  };

  const transitionMutationCallback = (mutationsList) => {
    for (const mutation of mutationsList) {
      const { type, attributeName, oldValue, target } = mutation;

      if (type !== "attributes" || attributeName !== ATTRIBUTE || !target) {
        return;
      }

      const state = target.getAttribute(ATTRIBUTE);
      addStateClasses(target, `data-transition-${state}`);

      if (state === "open") {
        showElement(target as HTMLElement);
      } else {
        const delay = getDelayFromAttribute(
          target,
          `data-transition-${oldValue}`
        );
        setTimeout(() => hideElement(target as HTMLElement), delay);
      }

      // Remove previous state classes
      if (oldValue) {
        removeStateClasses(target, `data-transition-${oldValue}`);
      }
    }
  };

  const observer = new MutationObserver(transitionMutationCallback);

  document.addEventListener("DOMContentLoaded", function () {
    const allElements = document.querySelectorAll(`[${ATTRIBUTE}]`);

    allElements.forEach((transitionElement: Element) => {
      const state = transitionElement.getAttribute(ATTRIBUTE);
      if (state !== "closed" && state !== "open") {
        throw new Error(
          `Invalid ${state} for transitions, use data-styles-[state] for non open/close transitions`
        );
      }
      if (state === "closed") {
        hideElement(transitionElement as HTMLElement);
      } else {
        showElement(transitionElement as HTMLElement);
      }
      addStateClasses(transitionElement, `data-transition-${state}`);
      observer.observe(transitionElement, {
        attributes: true,
        attributeOldValue: true,
        attributeFilter: [ATTRIBUTE],
      });
    });
  });
}
