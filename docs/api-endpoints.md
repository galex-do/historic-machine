# API Endpoints

Base URL: `http://localhost:8080/api`

All write endpoints require a valid JWT token in the `Authorization: Bearer <token>` header. Access level requirements are noted per endpoint.

---

## Authentication

| Method | Path | Description | Access |
|--------|------|-------------|--------|
| `POST` | `/auth/login` | Authenticate and receive a JWT token | Public |
| `POST` | `/auth/register` | Register a new account | Public |
| `POST` | `/auth/logout` | Invalidate the current session | Authenticated |
| `GET` | `/auth/me` | Get the current user's profile | Authenticated |

---

## Events

| Method | Path | Description | Access |
|--------|------|-------------|--------|
| `GET` | `/events` | List events with optional filtering (`locale`, `from`, `to`, `era`, `tags`, `lens_type`) | Public |
| `GET` | `/events/{id}` | Get a single event by ID | Public |
| `POST` | `/events` | Create a new event | User+ |
| `PUT` | `/events/{id}` | Update an event | Editor+ |
| `DELETE` | `/events/{id}` | Delete an event | Editor+ |
| `GET` | `/events/{id}/tags` | Get tags for an event | Public |
| `POST` | `/events/{id}/tags` | Set tags for an event (replaces existing) | Editor+ |

---

## Tags

| Method | Path | Description | Access |
|--------|------|-------------|--------|
| `GET` | `/tags` | List all tags | Public |
| `POST` | `/tags` | Create a new tag | Editor+ |
| `PUT` | `/tags/{id}` | Update a tag | Editor+ |
| `DELETE` | `/tags/{id}` | Delete a tag | Editor+ |

---

## Date Templates

| Method | Path | Description | Access |
|--------|------|-------------|--------|
| `GET` | `/date-template-groups` | List all template groups with their templates | Public |
| `GET` | `/date-templates` | List all templates | Public |
| `GET` | `/date-templates/{id}` | Get a single template | Public |
| `POST` | `/date-template-groups` | Create a template group | Editor+ |
| `PUT` | `/date-template-groups/{id}` | Update a template group | Editor+ |
| `DELETE` | `/date-template-groups/{id}` | Delete a template group (fails if templates exist) | Editor+ |
| `POST` | `/date-templates` | Create a template | Editor+ |
| `PUT` | `/date-templates/{id}` | Update a template | Editor+ |
| `DELETE` | `/date-templates/{id}` | Delete a template | Editor+ |

---

## Datasets

| Method | Path | Description | Access |
|--------|------|-------------|--------|
| `GET` | `/datasets` | List all datasets | Editor+ |
| `POST` | `/datasets/import` | Import a JSON dataset file | Editor+ |
| `PUT` | `/datasets/{id}/reset-modified` | Clear the modified flag after export | Editor+ |
| `DELETE` | `/datasets/{id}` | Delete a dataset and all its events | Editor+ |

---

## Regions

| Method | Path | Description | Access |
|--------|------|-------------|--------|
| `GET` | `/regions` | List all regions | Public |
| `GET` | `/regions/{id}` | Get a single region | Public |
| `POST` | `/regions` | Create a region | Editor+ |
| `PUT` | `/regions/{id}` | Update a region | Editor+ |
| `DELETE` | `/regions/{id}` | Delete a region | Editor+ |

---

## Users (Admin)

| Method | Path | Description | Access |
|--------|------|-------------|--------|
| `GET` | `/users` | List all users | Admin+ |
| `GET` | `/users/{id}` | Get a user | Admin+ |
| `PUT` | `/users/{id}` | Update a user (level, active status) | Admin+ |
| `DELETE` | `/users/{id}` | Deactivate a user | Admin+ |

---

## System

| Method | Path | Description | Access |
|--------|------|-------------|--------|
| `GET` | `/config` | Get public configuration (contact email, etc.) | Public |
| `GET` | `/support` | Get support/donation credentials | Public |
| `GET` | `/metrics` | Prometheus metrics endpoint | Public |
| `GET` | `/health` | Health check | Public |
