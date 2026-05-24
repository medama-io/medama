# AGENTS.md

## Role

You are working on `dashboard`, the React Router SPA / React TypeScript dashboard for Medama. Keep the user experience consistent with an operational analytics product: dense, readable, responsive, and focused on recurring dashboard workflows.

## Success Criteria

- Routes, loaders, actions, session handling, API errors, and cache-busting behavior remain compatible unless the task requires a change.
- UI changes follow existing component, hook, route, CSS module, and theme patterns.
- API calls use the local typed API client and generated OpenAPI types.
- The dashboard can still generate types and build through package mise tasks.

## Constraints

- Prefer existing local components, Mantine components, Radix primitives, hooks, and icons before adding dependencies.
- Use CSS modules for component-scoped styles and existing global/theme tokens for shared layout, colors, typography, and media queries.
- Preserve accessibility behavior for dialogs, icon buttons, labels, loading states, forms, and route errors.
- Keep browser-only code guarded where needed; tooling still has server/prerender boundaries even though the app runs in React Router SPA mode.
- Avoid broad visual redesigns while fixing narrow functional issues.

## API and Generated Types

- `app/api/types.d.ts` is generated from `../core/openapi.yaml`. Do not edit it by hand.
- When `core/openapi.yaml` changes, run `mise run generate` from `dashboard` or `mise run dashboard:generate` from the root.
- Prefer existing API helper modules over ad hoc `fetch` calls.
- Keep request/response typing aligned with generated OpenAPI types.

## Tooling

- Use Bun through package mise tasks. Run these from `dashboard`.
- `mise run dev`: generate API types and start the dev server.
- `mise run build`: generate API types and build the app.
- `mise run lint`: run Biome with writes enabled.
- `mise run lint:ci`: run typecheck, Biome in CI mode, and build.
- If dependencies change, update the root Bun lockfile through the repo workflow.

## Validation Defaults

- For UI or route behavior changes, run the smallest relevant check that covers the touched area, usually `mise run lint:ci` or a narrower typecheck/build command when appropriate.
- For visible UI changes, verify the affected route in a browser when practical.
- For frontend dependency, routing, styling, or component changes, capture a before/after browser baseline when practical. Prefer a production build with Vite preview for reviewable comparisons; use dev mode mainly for debugging.
- When browser verification needs API data, run the local core server with temporary databases and CORS for the dashboard origin.
- In final notes, name the compared routes, server mode, and any meaningful visual differences. Save screenshots when they make the comparison easier to review.
- For generated type changes, run the generation task and review the generated diff.
- For documentation-only changes, `git diff --check` is enough.

## Stop Rules

- Stop before changing route structure, API client semantics, cookie/session behavior, or embedded build flow unless the task explicitly calls for it.
- Do not commit `build` output unless a task explicitly requires tracked generated output.
