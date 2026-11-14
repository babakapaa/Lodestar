ALTER TABLE core.setup_state
	ADD COLUMN hardware_checks jsonb,
	ADD COLUMN hardware_at     timestamptz;
