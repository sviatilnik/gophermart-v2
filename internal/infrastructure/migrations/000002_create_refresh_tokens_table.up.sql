begin;
CREATE TABLE if not exists refresh_tokens (
    token      TEXT PRIMARY KEY,
    user_id    UUID NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX if not exists idx_refresh_tokens_token ON refresh_tokens(token);
commit;