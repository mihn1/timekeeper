# Frontend (Svelte + Vite)

This directory contains the Wails frontend for TimeKeeper.

- Main app shell: `src/App.svelte`
- Primary views: `src/components/Dashboard.svelte`, `src/components/rules/Rules.svelte`, `src/components/categories/Categories.svelte`
- Shared UI: `src/components/common/`
- Stores: `src/stores/`
- Theme tokens: `src/theme.css`

For backend architecture and full project context, read `../../AGENTS.md`.

## Install and Run

From this directory:

```bash
npm install
npm run dev
```

Build production frontend bundle:

```bash
npm run build
```

## Integration With Wails

- Backend method bindings are generated into `../wailsjs/`.
- Use methods from `../wailsjs/go/main/App` in Svelte components/stores.
- Do not hand-edit files under `../wailsjs/`; regenerate via Wails flows.

## Conventions in This Project

- `refreshData` store and `timekeeper:data-updated` event are used to trigger view refresh.
- Rule/category CRUD screens rely on generated DTO types in `../wailsjs/go/models`.
- Theme is controlled via `data-theme` on `document.documentElement` and CSS variables in `src/theme.css`.

## Contributor / Agent Rules

- Keep generated bindings untouched; change Go API signatures in `ui/` then regenerate.
- Prefer aligning UI changes with existing component and store patterns.
- Do not create commits unless explicitly requested by the user.
