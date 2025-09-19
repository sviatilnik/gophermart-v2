begin;
CREATE TABLE wallet_events (
      event_id      TEXT PRIMARY KEY,
      aggregate_id  TEXT NOT NULL,
      event_type    TEXT NOT NULL,
      event_data    JSONB NOT NULL,
      version       INTEGER NOT NULL,
      timestamp     TIMESTAMP WITH TIME ZONE NOT NULL,

      CONSTRAINT version_sequence UNIQUE (aggregate_id, version)
);

CREATE INDEX idx_wallet_events_aggregate ON wallet_events (aggregate_id);
CREATE INDEX idx_wallet_events_version ON wallet_events (aggregate_id, version);
commit;