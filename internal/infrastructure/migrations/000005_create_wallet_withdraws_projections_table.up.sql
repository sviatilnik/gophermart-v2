CREATE TABLE wallet_withdrawals (
    id           TEXT PRIMARY KEY,
    event_id    TEXT NOT NULL,
    customer_id  TEXT NOT NULL,
    amount       FLOAT NOT NULL,
    order_number text not null,
    timestamp    TIMESTAMP WITH TIME ZONE NOT NULL,

    FOREIGN KEY (event_id) REFERENCES wallet_events(event_id) ON DELETE NO ACTION
);

CREATE INDEX idx_wallet_withdrawals_customer ON wallet_withdrawals(customer_id);