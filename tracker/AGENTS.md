# AGENTS.md

## Role

You are working on `tracker`, the tiny browser JavaScript tracker used by Medama sites. Preserve browser behavior, privacy properties, and compressed transfer size while making the smallest useful change.

## Success Criteria

- `src/tracker.js` remains the source of truth for runtime tracker behavior.
- Generated `dist` artifacts are updated only by the tracker build.
- Payload keys, event names, data attributes, and endpoint paths remain compatible with `core`.
- The default/minimal tracker stays under the documented gzip size target when possible.
- Runtime behavior is validated in real browsers for behavior changes.

## Constraints

- Keep runtime code dependency-free.
- Preserve the no-cookie tracker design and avoid adding identifiers beyond existing beacon/session mechanics.
- Prefer size-aware direct code over abstractions or helper layers.
- Measure compressed output before accepting a size-motivated change; smaller source or minified output can still compress worse.
- Keep feature-preprocessor conventions and alphabetic combined-feature filenames.

## Build Outputs

- `dist/*.js` and `dist/*.min.js` are generated tracked artifacts.
- `core/client/scripts/*` is embed output copied by `mise run embed` and is not the source of truth.
- Change `src/tracker.js` or build scripts first, run the build, then review generated output. Do not hand-edit `dist`.

## Tooling

- Run these from `tracker`.
- `mise run build`: rebuild tracker variants and print compressed sizes.
- `mise run test`: build tracker variants needed by fixtures and run Playwright tests.
- `mise run test:setup`: install Playwright browsers and dependencies when needed.
- `mise run format`: run Biome formatting.
- `mise run serve`: run the local fixture server for manual debugging.

## Validation Defaults

- Use Playwright for load/unload beacons, history navigation, hash mode, `data-api`, and custom click/page events.
- Keep fixtures focused on observable browser/API behavior.
- If a change affects tracker size, include the relevant `mise run build` size output in the final response.
- If a change affects API payloads, validate matching `core` endpoint behavior or update the API contract intentionally.

## Stop Rules

- Stop before changing payload keys, endpoint paths, privacy behavior, or runtime dependencies unless the task explicitly calls for it.
- Stop if browser validation is required but Playwright browsers or local services are unavailable; report the blocker and the next best check.
