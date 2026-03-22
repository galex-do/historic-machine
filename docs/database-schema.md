# Database Schema

timediverr uses PostgreSQL. Migrations are managed with [Goose](https://github.com/pressly/goose) and live in `backend/migrations/`.

---

## Tables

### `events`
Core table for historical events.

| Column | Type | Notes |
|--------|------|-------|
| `id` | `SERIAL PK` | |
| `name` | `VARCHAR(255)` | English name (default display) |
| `name_ru` | `VARCHAR(255)` | Russian name |
| `description` | `TEXT` | English description |
| `description_ru` | `TEXT` | Russian description |
| `source` | `TEXT` | Source URL or reference |
| `latitude` | `DECIMAL(10,8)` | |
| `longitude` | `DECIMAL(11,8)` | |
| `event_date` | `DATE` | Stored as PostgreSQL DATE |
| `era` | `VARCHAR(2)` | `'BC'` or `'AD'` |
| `lens_type` | `VARCHAR(50)` | Category: `historic`, `political`, `cultural`, `military`, `scientific`, `religious` |
| `dataset_id` | `INTEGER FK → event_datasets` | Nullable |
| `created_by` | `INTEGER FK → users` | Nullable |
| `updated_by` | `INTEGER FK → users` | Nullable |
| `created_at` | `TIMESTAMP` | |
| `updated_at` | `TIMESTAMP` | |

---

### `tags`
Flexible tagging system with visual and behavioural options.

| Column | Type | Notes |
|--------|------|-------|
| `id` | `SERIAL PK` | |
| `name` | `VARCHAR(100) UNIQUE` | |
| `description` | `TEXT` | |
| `color` | `VARCHAR(7)` | Hex colour, default `#3B82F6` |
| `border_color` | `VARCHAR(7)` | Optional inner border via `box-shadow: inset` |
| `key_color` | `BOOLEAN` | When true, shows a coloured dot next to event names in timeline/cluster views |
| `emoji` | `VARCHAR(10)` | When set, overrides the default lens-type emoji on map markers |
| `weight` | `INTEGER` | Ordering weight; higher = shown first |
| `created_at` | `TIMESTAMP` | |
| `updated_at` | `TIMESTAMP` | |

---

### `event_tags`
Many-to-many junction between events and tags.

| Column | Type | Notes |
|--------|------|-------|
| `event_id` | `INTEGER FK → events` | Composite PK |
| `tag_id` | `INTEGER FK → tags` | Composite PK |

---

### `users`

| Column | Type | Notes |
|--------|------|-------|
| `id` | `SERIAL PK` | |
| `username` | `VARCHAR(100) UNIQUE` | |
| `email` | `VARCHAR(255)` | Optional |
| `password_hash` | `VARCHAR(255)` | bcrypt |
| `access_level` | `VARCHAR(20)` | `guest` / `user` / `editor` / `admin` / `super` |
| `is_active` | `BOOLEAN` | |
| `created_at` | `TIMESTAMP` | |
| `updated_at` | `TIMESTAMP` | |
| `last_login` | `TIMESTAMP` | |

---

### `user_sessions`
Active JWT sessions for authenticated users.

| Column | Type | Notes |
|--------|------|-------|
| `id` | `SERIAL PK` | |
| `user_id` | `INTEGER FK → users` | |
| `token_hash` | `VARCHAR(255) UNIQUE` | |
| `expires_at` | `TIMESTAMP` | |
| `last_seen_at` | `TIMESTAMP` | Updated on each heartbeat (60s interval) |
| `is_active` | `BOOLEAN` | |
| `created_at` | `TIMESTAMP` | |

---

### `event_datasets`
Tracks imported event collections.

| Column | Type | Notes |
|--------|------|-------|
| `id` | `SERIAL PK` | |
| `filename` | `VARCHAR(255)` | Original import filename |
| `description` | `TEXT` | |
| `event_count` | `INTEGER` | Set at import time |
| `modified` | `BOOLEAN` | Set to `true` automatically when any associated event is created, updated, or deleted. Reset via API after re-export. |
| `uploaded_by` | `INTEGER FK → users` | Nullable |
| `created_at` | `TIMESTAMP` | |
| `updated_at` | `TIMESTAMP` | |

---

### `date_template_groups`
Named groups of date range templates (e.g. "Ancient Greece", "Roman Empire").

| Column | Type | Notes |
|--------|------|-------|
| `id` | `SERIAL PK` | |
| `name` | `VARCHAR(255) UNIQUE` | English name |
| `name_ru` | `VARCHAR(255)` | Russian name |
| `description` | `TEXT` | |
| `description_ru` | `TEXT` | |
| `display_order` | `INTEGER` | Sort order in the selector |
| `created_at` | `TIMESTAMP` | |

---

### `date_templates`
Individual date range presets within a group.

| Column | Type | Notes |
|--------|------|-------|
| `id` | `SERIAL PK` | |
| `group_id` | `INTEGER FK → date_template_groups` | Cascade delete |
| `name` | `VARCHAR(255)` | English name |
| `name_ru` | `VARCHAR(255)` | Russian name |
| `description` | `TEXT` | |
| `description_ru` | `TEXT` | |
| `start_date` | `DATE` | |
| `start_era` | `VARCHAR(2)` | `'BC'` or `'AD'` |
| `end_date` | `DATE` | |
| `end_era` | `VARCHAR(2)` | `'BC'` or `'AD'` |
| `display_order` | `INTEGER` | Sort order within group |
| `created_at` | `TIMESTAMP` | |

---

### `regions`
Polygonal map overlays tied to date templates (e.g. empires, territories).

| Column | Type | Notes |
|--------|------|-------|
| `id` | `SERIAL PK` | |
| `name` | `VARCHAR(255)` | |
| `name_en` | `VARCHAR(255)` | |
| `name_ru` | `VARCHAR(255)` | |
| `description` | `TEXT` | |
| `geojson` | `JSONB` | GeoJSON polygon geometry |
| `color` | `VARCHAR(7)` | Fill colour |
| `fill_opacity` | `REAL` | Default `0.2` |
| `border_color` | `VARCHAR(7)` | |
| `border_width` | `REAL` | Default `2` |
| `created_at` | `TIMESTAMP` | |
| `updated_at` | `TIMESTAMP` | |

---

### `template_regions`
Many-to-many: which regions appear when a template is active.

| Column | Type | Notes |
|--------|------|-------|
| `template_id` | `INTEGER FK → date_templates` | Composite PK |
| `region_id` | `INTEGER FK → regions` | Composite PK |

---

### `support_credentials`
Configurable donation/support links shown on the About page.

| Column | Type | Notes |
|--------|------|-------|
| `id` | `SERIAL PK` | |
| `name` | `VARCHAR(255)` | Display label |
| `value` | `TEXT` | URL or address |
| `is_url` | `BOOLEAN` | If true, rendered as a clickable link |
| `created_at` | `TIMESTAMP` | |

---

## Views

### `events_with_display_dates`
Extends `events` with:
- `display_date` — human-readable date string including era (e.g. `"21.04.0753 BC"`)
- `astronomical_year` — signed integer for correct BC/AD sorting
- `tags` — aggregated JSON array of all tags with `id`, `name`, `color`, `border_color`, `key_color`, `emoji`, `weight`

### `date_templates_with_display`
Extends `date_templates` with:
- `start_display_date` / `end_display_date` — formatted date strings
- `start_astronomical_year` / `end_astronomical_year` — for sorting
- `group_name` — joined from `date_template_groups`

---

## Running Migrations

```bash
# Via Docker Compose
docker compose exec backend goose -dir migrations postgres \
  "postgres://postgres:password@db:5432/historical_events?sslmode=disable" up

# Or with Make
make migrate
```
