import initAdminNav from "./admin/admin.nav";
import initNav from "./app/nav";
import initMycodes from "./app/myCodes";
import initToast from "./shared/toast";
import initTailwindTransitions from "./tailwind/tailwind-transitions";
import initCatalogFilters from "./app/catalogFilters";

initTailwindTransitions();
initMycodes();
initAdminNav();
initNav();
initToast();
initCatalogFilters();
