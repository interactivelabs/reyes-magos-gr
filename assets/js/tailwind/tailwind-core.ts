export enum Action {
  ADD = "add",
  REMOVE = "remove",
}

export const getStateClasses = (element: Element, attribute: string) => {
  const stateClasses = element.getAttribute(attribute);
  if (!stateClasses) {
    throw new Error(`Missing ${attribute} attribute`);
  }
  const stateClassesArray = stateClasses.split(" ");
  const afterClasses = stateClassesArray.filter((className) =>
    className.startsWith("final:")
  );
  const immediateClasses = stateClassesArray.filter(
    (className) => !className.startsWith("final:")
  );
  return [immediateClasses, afterClasses];
};

export const updateStateClasses = (
  element: Element,
  attribute: string,
  action: Action
) => {
  const [immediateClasses, afterClasses] = getStateClasses(element, attribute);

  const createTimeoutToAddClass = (delay: number, className: string) => {
    setTimeout(() => {
      element.classList.add(className);
    }, delay);
  };

  immediateClasses.forEach((className) => {
    if (action === Action.ADD) {
      element.classList.add(className);
      createTimeoutToAddClass(10, className);
    } else {
      element.classList.remove(className);
    }
  });

  afterClasses.forEach((className) => {
    const props = className.split("final:")[1];
    const [delay, afterClassName] = JSON.parse(props.replace(/'/g, '"'));
    if (action === Action.ADD) {
      createTimeoutToAddClass(delay, afterClassName);
    } else {
      element.classList.remove(afterClassName);
    }
  });
};

export const addStateClasses = (element: Element, attribute: string) => {
  updateStateClasses(element, attribute, Action.ADD);
};

export const removeStateClasses = (element: Element, attribute: string) => {
  updateStateClasses(element, attribute, Action.REMOVE);
};

export const getDelayFromAttribute = (element: Element, attribute: string) => {
  const classNames = element.getAttribute(attribute);
  if (!classNames) {
    throw new Error(`Missing ${attribute} attribute`);
  }
  const classNamesArray = classNames.split(" ");
  const delayClasses = classNamesArray.filter(
    (className) =>
      className.startsWith("delay-") || className.startsWith("duration-")
  );

  if (!delayClasses || !delayClasses.length) return 150;

  let delay = 0;

  delayClasses.forEach((delayClass) => {
    const delayValue = parseInt(delayClass.split("-")[1], 10);
    delay += delayValue;
  });

  return delay;
};
