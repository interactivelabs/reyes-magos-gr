/// <reference types="alpinejs" />
import type { Alpine as AlpineType } from "alpinejs";

declare global {
  var Alpine: AlpineType;
}

import "./shared/htmxErrorHandler";
import "./shared/toasts";

import "./app/myCodes";
import "./app/catalog";
import "./app/homeAnnimations";
import "./app/redeem";
