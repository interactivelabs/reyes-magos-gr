/// <reference types="alpinejs" />
import type { Alpine as AlpineType } from "alpinejs";

declare global {
  var Alpine: AlpineType;
}

import "./components/toasts";
import "./components/search-box";

import "./shared/htmxErrorHandler";

import "./app/myCodes";
import "./app/catalog";
import "./app/homeAnnimations";
import "./app/redeem";
