# Alpine.js Deprecation List

Goal: remove Alpine.js in favor of Tailwind, htmx, and vanilla TS.

## Script includes

- `views/layout.templ` (lines 41-45) — loads Alpine core + `@alpinejs/focus` + `@alpinejs/collapse` CDN scripts. Remove once all usages below are migrated.

## Files using Alpine.js

| File | Alpine features used |
|---|---|
| `views/shared/Nav.templ` | `x-data`, `x-show`, `x-cloak`, `x-transition`, `@click`, `@click.away` — mobile menu toggle + admin dropdown |
| `views/components/toast.templ` | Heaviest user: `x-data="toasts"` (global Alpine.data component), `x-for`, `x-if`, `x-show`, `x-cloak`, `x-transition`, `x-init`, `x-trap` (focus plugin), `$dispatch`, `$nextTick`, `x-text` |
| `views/checkout/checkout.templ` | `x-data="codeForm"` (Alpine.data component) + `@submit`, `x-show` for form validation; separate `x-data` per toy-image-selector with `x-bind:class`, `@click` |
| `views/checkout/accordion.templ` | `x-data`, `x-show`, `x-collapse` (collapse plugin), `x-cloak`, `@click` |
| `views/admin/toys/autocomplete_categories.templ` | `x-data="CategoryWithSearch(...)"` (Alpine.data component), `x-model`, `x-show`, `x-for`, `@click` |
| `views/admin/toys/toys.templ`, `create_toy_form.templ`, `update_toy_form.templ` | `x-on:htmx:after-swap.window="modalOpen = true"` / `x-on:htmx:after-on-load="modalOpen = false"` |
| `views/admin/volunteers/volunteers.templ`, `create_volunteer_form.templ`, `update_volunteer_form.templ` | Same `modalOpen` pattern as above |

## Backing JS (Alpine.data components — will also need rewriting)

- `assets/js/components/toasts.ts` → `Alpine.data("toasts", ...)`
- `assets/js/app/checkout.ts` → `Alpine.data("codeForm", ...)`
- `assets/js/admin/toys_form.ts` → `Alpine.data("CategoryWithSearch", ...)`

## Notable oddity

`modalOpen` (used in the toys/volunteers admin modal templates) is referenced via `x-on:...="modalOpen = true/false"` but there's no `x-data` ancestor declaring `modalOpen` anywhere in the views. It's relying on Alpine implicitly creating a loose global — worth double-checking this actually does anything today, since it may already be dead/no-op code from a half-finished modal implementation.

## Migration difficulty ranking

1. **Easy** (pure show/hide toggles, replaceable with a `<details>`/`<dialog>` element or a tiny vanilla-TS class toggle): `Nav.templ`, `catalog.templ` back-to-top, `accordion.templ`
2. **Medium** (small local state + htmx events): the `modalOpen` admin modal triggers, checkout image-selector `x-data`
3. **Harder** (full component logic - search/autocomplete, dynamic list rendering, toast queue with timers/animations): `toast.templ` + `toasts.ts`, `autocomplete_categories.templ` + `toys_form.ts`, `checkout.templ` codeForm + `checkout.ts`
