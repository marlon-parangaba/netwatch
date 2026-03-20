-- ============================================================
-- NetWatch - Migration 001: Schema inicial
-- ============================================================

-- Habilitar extensões necessárias
CREATE EXTENSION IF NOT EXISTS "pgcrypto";    -- para gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "pg_trgm";     -- para busca por similaridade

-- ============================================================
-- USERS
-- ============================================================
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(20)  NOT NULL DEFAULT 'viewer'
                  CHECK (role IN ('admin', 'operator', 'viewer')),
    active        BOOLEAN NOT NULL DEFAULT true,
    last_login_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_users_email ON users (email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_deleted_at ON users (deleted_at);

-- ============================================================
-- DEVICES
-- ============================================================
CREATE TABLE devices (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name             VARCHAR(255) NOT NULL,
    hostname         VARCHAR(255),
    ip_address       VARCHAR(45) NOT NULL,
    type             VARCHAR(50) NOT NULL DEFAULT 'generic'
                     CHECK (type IN ('mikrotik','cisco','huawei','juniper','ubiquiti','generic')),
    status           VARCHAR(20) NOT NULL DEFAULT 'unknown'
                     CHECK (status IN ('online','offline','warning','unknown')),
    description      TEXT,
    location         VARCHAR(255),

    -- SNMP Config
    snmp_version     VARCHAR(5) NOT NULL DEFAULT 'v2c'
                     CHECK (snmp_version IN ('v1','v2c','v3')),
    snmp_community   VARCHAR(100) DEFAULT 'public',
    snmp_port        INTEGER NOT NULL DEFAULT 161,

    -- SNMPv3
    snmpv3_username  VARCHAR(100),
    snmpv3_auth_proto VARCHAR(10),
    snmpv3_auth_key  VARCHAR(255),
    snmpv3_priv_proto VARCHAR(10),
    snmpv3_priv_key  VARCHAR(255),

    -- Polling config
    polling_interval INTEGER NOT NULL DEFAULT 60,
    enabled          BOOLEAN NOT NULL DEFAULT true,

    -- Sys info (from SNMP)
    sys_name         VARCHAR(255),
    sys_descr        TEXT,
    sys_oid          VARCHAR(255),
    sys_contact      VARCHAR(255),

    -- Tags customizadas
    tags             JSONB NOT NULL DEFAULT '[]',

    -- Status timestamps
    last_seen_at     TIMESTAMPTZ,
    last_polled_at   TIMESTAMPTZ,

    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_devices_ip ON devices (ip_address) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_status ON devices (status) WHERE deleted_at IS NULL;
CREATE INDEX idx_devices_deleted_at ON devices (deleted_at);

-- ============================================================
-- INTERFACE INFO
-- ============================================================
CREATE TABLE interface_infos (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id        UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    if_index         INTEGER NOT NULL,
    if_descr         VARCHAR(255),
    if_type          INTEGER,
    if_speed         BIGINT,
    if_admin_status  INTEGER,
    if_oper_status   INTEGER,
    if_alias         VARCHAR(255),
    ip_address       VARCHAR(45),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (device_id, if_index)
);

CREATE INDEX idx_interfaces_device ON interface_infos (device_id);

-- ============================================================
-- METRICS (Particionado por mês - retenção 1 ano)
-- ============================================================
CREATE TABLE metrics (
    id             UUID NOT NULL DEFAULT gen_random_uuid(),
    device_id      UUID NOT NULL,
    type           VARCHAR(50) NOT NULL,
    interface_index INTEGER NOT NULL DEFAULT 0,
    oid            VARCHAR(255),
    value          DOUBLE PRECISION NOT NULL,
    unit           VARCHAR(20),
    collected_at   TIMESTAMPTZ NOT NULL
) PARTITION BY RANGE (collected_at);

-- Índices na tabela pai
CREATE INDEX idx_metrics_device_type ON metrics (device_id, type, collected_at DESC);
CREATE INDEX idx_metrics_collected_at ON metrics (collected_at DESC);

-- Função para criar partições de métricas automaticamente
CREATE OR REPLACE FUNCTION create_metrics_partition(target_date DATE)
RETURNS VOID AS $$
DECLARE
    partition_name TEXT;
    start_date DATE;
    end_date DATE;
BEGIN
    partition_name := 'metrics_' || TO_CHAR(target_date, 'YYYY_MM');
    start_date := DATE_TRUNC('month', target_date)::DATE;
    end_date := (DATE_TRUNC('month', target_date) + INTERVAL '1 month')::DATE;

    IF NOT EXISTS (
        SELECT 1 FROM pg_tables WHERE tablename = partition_name
    ) THEN
        EXECUTE FORMAT(
            'CREATE TABLE %I PARTITION OF metrics FOR VALUES FROM (%L) TO (%L)',
            partition_name, start_date, end_date
        );
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Criar partições para os próximos 13 meses
DO $$
DECLARE i INTEGER;
BEGIN
    FOR i IN 0..12 LOOP
        PERFORM create_metrics_partition(
            (DATE_TRUNC('month', NOW()) + (i || ' months')::INTERVAL)::DATE
        );
    END LOOP;
END $$;

-- ============================================================
-- ALERT RULES
-- ============================================================
CREATE TABLE alert_rules (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name           VARCHAR(255) NOT NULL,
    description    TEXT,
    device_id      UUID REFERENCES devices(id) ON DELETE CASCADE,
    metric_type    VARCHAR(50) NOT NULL,
    condition      VARCHAR(10) NOT NULL
                   CHECK (condition IN ('gt','gte','lt','lte','eq','neq','down')),
    threshold      DOUBLE PRECISION NOT NULL DEFAULT 0,
    severity       VARCHAR(20) NOT NULL DEFAULT 'medium'
                   CHECK (severity IN ('critical','high','medium','low','info')),
    enabled        BOOLEAN NOT NULL DEFAULT true,
    notifications  JSONB NOT NULL DEFAULT '[]',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ
);

CREATE INDEX idx_alert_rules_device ON alert_rules (device_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_alert_rules_deleted_at ON alert_rules (deleted_at);

-- ============================================================
-- ALERT EVENTS (Retenção: 6 meses rolling window)
-- ============================================================
CREATE TABLE alert_events (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_id      UUID NOT NULL REFERENCES alert_rules(id) ON DELETE CASCADE,
    device_id    UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    status       VARCHAR(20) NOT NULL DEFAULT 'active'
                 CHECK (status IN ('active','resolved','acknowledged')),
    severity     VARCHAR(20) NOT NULL,
    value        DOUBLE PRECISION,
    message      TEXT,
    triggered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at  TIMESTAMPTZ,
    acked_at     TIMESTAMPTZ,
    acked_by_id  UUID REFERENCES users(id),
    duration_seconds BIGINT
);

CREATE INDEX idx_alert_events_device ON alert_events (device_id, triggered_at DESC);
CREATE INDEX idx_alert_events_rule ON alert_events (rule_id);
CREATE INDEX idx_alert_events_status ON alert_events (status) WHERE status = 'active';
CREATE INDEX idx_alert_events_triggered ON alert_events (triggered_at DESC);

-- ============================================================
-- DOWNTIME EVENTS (Retenção: vitalícia - para MTBF/MTTR)
-- ============================================================
CREATE TABLE downtime_events (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id      UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    started_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at       TIMESTAMPTZ,
    duration_seconds BIGINT,
    cause          VARCHAR(255),
    notes          TEXT
);

CREATE INDEX idx_downtime_device ON downtime_events (device_id, started_at DESC);
CREATE INDEX idx_downtime_active ON downtime_events (device_id) WHERE ended_at IS NULL;

-- ============================================================
-- NETWORK MAPS
-- ============================================================
CREATE TABLE network_maps (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    layout      JSONB NOT NULL DEFAULT '{"nodes":[],"edges":[]}',
    created_by  UUID REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_maps_deleted_at ON network_maps (deleted_at);

-- ============================================================
-- Função para atualizar updated_at automaticamente
-- ============================================================
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_devices_updated_at BEFORE UPDATE ON devices
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_alert_rules_updated_at BEFORE UPDATE ON alert_rules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
CREATE TRIGGER trg_maps_updated_at BEFORE UPDATE ON network_maps
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
