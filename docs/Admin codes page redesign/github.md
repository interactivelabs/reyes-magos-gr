repo: interactivelabs/reyes-magos-gr
branch: main

## Last sync
date: 2026-08-02T23:16:13Z

### Updated in this project
- Redesigned the admin "Codes" page (behind auth) as Admin Codes.dc.html
- Applied the in-repo "soft & hopeful" design spec (docs/DL Toys Design Spec.dc.html) — peach/lavender palette, Quicksand + Cabin type, shadow-pop buttons, mono code chips, 20px card radius
- UX changes: always-visible stats (unassigned/assigned/total), create-codes always available (no longer hidden when unassigned codes exist), header checkbox "select all", live selected-count on Assign/Quitar buttons, empty states for both tables, inline success messages
- Added tweaks: create-codes placement (inline card vs modal) and table density (comfortable vs compact)

## Screen map
| Screen | Repo source |
|---|---|
| Admin Codes.dc.html | views/admin/codes/codes_templ.go, handlers/admin/codes.go, store/codes_store.go, store/models/codes.go, services/codes.go, views/admin/admin_layout_templ.go, docs/DL Toys Design Spec.dc.html |
