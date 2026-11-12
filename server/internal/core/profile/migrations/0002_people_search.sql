CREATE INDEX users_name_search ON core.users USING gin (core.search_norm(display_name) gin_trgm_ops);
CREATE INDEX profiles_name_search ON core.profiles USING gin (core.search_norm(display_name) gin_trgm_ops);
