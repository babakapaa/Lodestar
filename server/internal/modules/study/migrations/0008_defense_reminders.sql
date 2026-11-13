ALTER TABLE study.defenses ADD COLUMN reminded_slot uuid;

CREATE INDEX defense_slots_start ON study.defense_slots (tenant_id, starts_at);
