# User Access Levels

timediverr uses role-based access control with four levels. Each level inherits all permissions of the levels below it.

| Level | Description |
|-------|-------------|
| `guest` | Read-only access. Can browse the map, view events, filter by tags and date ranges, and open event details. No account required. |
| `user` | All guest permissions plus the ability to create new events and add tags. Cannot edit or delete events created by others. |
| `editor` | All user permissions plus full event management (edit/delete any event), dataset import/export, and access to the admin panel for tags and templates. |
| `admin` | All editor permissions plus user account management (create, edit, deactivate users). Cannot promote users to `super`. |
| `super` | Full system access including all admin capabilities and the ability to manage `super`-level accounts, configure support credentials, and access system metrics. |

## Creating the First Super User

The first super user must be created directly via the database since there is no super user to grant it through the UI:

```bash
docker compose exec db psql -U postgres -d historical_events -c \
  "INSERT INTO users (username, password_hash, access_level) VALUES
   ('admin', crypt('your_password', gen_salt('bf', 12)), 'super');"
```

Subsequent users can be registered normally and promoted through the admin panel.
