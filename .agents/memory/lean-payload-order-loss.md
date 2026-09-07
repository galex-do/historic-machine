---
name: Lean payload refactors silently drop implied ordering
description: When a bulk/map API payload is slimmed from full embedded objects down to bare IDs, any ordering that used to ride along with the embedded objects (e.g. tags sorted by weight) must be explicitly re-derived by the resolver — it does not survive the ID roundtrip.
---

## The lesson

The events list endpoint (`backend/internal/models/event.go`, `EventListItem`) was refactored to send only `tag_ids` per event instead of full tag objects, to cut map payload size. The backend still builds `tag_ids` in weight-sorted order (from the `events_with_display_dates` view's `JSON_AGG(... ORDER BY t.weight DESC, t.name)`), but that ordering is an implicit, undocumented contract — nothing enforces it downstream.

On the frontend, `getTagsByIds` (`frontend/src/composables/useTags.js`) resolved `tag_ids` back into full tag objects via `allTags.value.filter(tag => ids.includes(tag.id))`. Since `allTags` itself is cached from `GET /api/tags`, which is sorted alphabetically by name (`ORDER BY name ASC` in `tag_repository.go`), the filter silently returned tags in name order, not weight order — breaking anything that assumed "first tag = highest weight" (e.g. `getTagEmoji` in `event-utils.js`).

**Why:** A payload-shrinking refactor that reduces embedded objects to bare IDs looks lossless (the resolver can always look the ID up), but any ordering guarantee tied to the original embedding is lost unless the resolver explicitly re-establishes it. This is easy to miss because nothing errors — the resolved objects are correct, just reordered.

**How to apply:** When reviewing or writing a resolver that turns bulk IDs back into full objects from a separately-cached catalog, check whether call sites depend on result order (first-element assumptions, display priority, sort-by-weight/priority patterns). If they do, sort explicitly in the resolver using the authoritative field (e.g. `weight`) rather than trusting either the ID array's order or the cached catalog's fetch order.
