# Alpine.js Deprecation List

Goal: remove Alpine.js in favor of Tailwind, htmx, and vanilla TS.

## Script includes

- `views/layout.templ` (lines 41-45) — loads Alpine core + `@alpinejs/focus` + `@alpinejs/collapse` CDN scripts. Remove once all usages below are migrated.

## Files using Alpine.js

| File | Alpine features used |
|---|---|
| `views/components/toast.templ` | Heaviest user: `x-data="toasts"` (global Alpine.data component), `x-for`, `x-if`, `x-show`, `x-cloak`, `x-transition`, `x-init`, `x-trap` (focus plugin), `$dispatch`, `$nextTick`, `x-text` |

## Backing JS (Alpine.data components — will also need rewriting)

- `assets/js/components/toasts.ts` → `Alpine.data("toasts", ...)`

## Notable oddity

`modalOpen` (used in the toys/volunteers admin modal templates) is referenced via `x-on:...="modalOpen = true/false"` but there's no `x-data` ancestor declaring `modalOpen` anywhere in the views. It's relying on Alpine implicitly creating a loose global — worth double-checking this actually does anything today, since it may already be dead/no-op code from a half-finished modal implementation.

## Migration difficulty ranking

1. **Easy** (pure show/hide toggles, replaceable with a `<details>`/`<dialog>` element or a tiny vanilla-TS class toggle): `Nav.templ`, `catalog.templ` back-to-top, `accordion.templ`
2. **Medium** (small local state + htmx events): the `modalOpen` admin modal triggers, checkout image-selector `x-data`
3. **Harder** (full component logic - search/autocomplete, dynamic list rendering, toast queue with timers/animations): `toast.templ` + `toasts.ts`
