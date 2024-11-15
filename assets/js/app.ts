/// <reference types="alpinejs" />
import type { Alpine as AlpineType } from "alpinejs";

declare global {
  var Alpine: AlpineType;
}

import "./app/myCodes";
import "./shared/toasts";
import "./app/catalogFilters";
import "./shared/htmxErrorHandler";
