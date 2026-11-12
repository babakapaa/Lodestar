CREATE SCHEMA IF NOT EXISTS core;
GRANT USAGE ON SCHEMA core TO lodestar_app;

CREATE TABLE core.tenants (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	slug       text        NOT NULL UNIQUE CHECK (slug ~ '^[a-z0-9][a-z0-9-]{1,39}$'),
	name       text        NOT NULL CHECK (length(name) BETWEEN 1 AND 200),
	created_at timestamptz NOT NULL DEFAULT now()
);

GRANT SELECT ON core.tenants TO lodestar_app;
