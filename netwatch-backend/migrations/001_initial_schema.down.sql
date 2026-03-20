-- Rollback da migration 001
DROP TRIGGER IF EXISTS trg_maps_updated_at ON network_maps;
DROP TRIGGER IF EXISTS trg_alert_rules_updated_at ON alert_rules;
DROP TRIGGER IF EXISTS trg_devices_updated_at ON devices;
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at();
DROP FUNCTION IF EXISTS create_metrics_partition(DATE);
DROP TABLE IF EXISTS network_maps;
DROP TABLE IF EXISTS downtime_events;
DROP TABLE IF EXISTS alert_events;
DROP TABLE IF EXISTS alert_rules;
DROP TABLE IF EXISTS metrics;
DROP TABLE IF EXISTS interface_infos;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS pg_trgm;
DROP EXTENSION IF EXISTS pgcrypto;
