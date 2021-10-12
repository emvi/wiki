# Emvi Wiki

**Support?**

No. This used to be our SaaS on emvi.com, but didn't work out the way we wanted. It's now open-source to help some of our users keep it running, but support will be *very* limited.

## Hosts configuration

In order for Emvi to work locally, the hosts file (`/etc/hosts`) must be manipulated. Add the following line to the end of the file:

```
127.0.0.1   localhost.com dev.localhost.com dev1.localhost.com dev2.localhost.com dev3.localhost.com
```

Emvi can then be accessed on `localhost.com`.

## Database Setup

Make sure you run the `backend/schema/manually.sql` to setup the backend database.

## Development Database Data

The following data must exist in the auth database in order for the local development environment to work:

```
INSERT INTO "client" (name, client_id, client_secret, redirect_uri, trusted) VALUES
('backend', '9lrodHgg0z4EGCe9dTJ1', 'lFONmmtVeHoh7yLUzKXWGG4WnfUCMWNejrrxWlm8ZDSp5Sjlu8PVNSvsGm0ju30d', 'http://localhost.com:4000/organizations', true),
('collab', 'xz5tN33UW6kZzIyHrO8x', '3DLvZUadCoy6xB9mHx5uaJkfsY3K7WjwVBW6pdAeTAEy37Bjf4u05By7e0QT0hBT', null, true);
```
