import initMycodes from "./app/myCodes";
import initToast from "./shared/toast";
import initTailwindTransitions from "./tailwind/tailwind-transitions";
import initCatalogFilters from "./app/catalogFilters";
import "./shared/htmxErrorHandler";

initTailwindTransitions();
initMycodes();
initToast();
initCatalogFilters();
