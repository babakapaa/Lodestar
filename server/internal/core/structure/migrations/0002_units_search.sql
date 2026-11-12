CREATE INDEX units_name_search ON core.units USING gin (core.search_norm(name) gin_trgm_ops) WHERE archived_at IS NULL;
CREATE INDEX units_code_search ON core.units USING gin (core.search_norm(coalesce(code, '')) gin_trgm_ops) WHERE archived_at IS NULL;
