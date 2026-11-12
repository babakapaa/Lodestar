CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE FUNCTION core.search_norm(value text) RETURNS text
	LANGUAGE sql IMMUTABLE PARALLEL SAFE STRICT
	RETURN lower(translate(value, 'Ёё', 'Ее'));

GRANT EXECUTE ON FUNCTION core.search_norm(text) TO lodestar_app;
