create table if not exists orders(
    id uuid primary key,
    number text not null unique,
    user_id uuid not null,
    state varchar(255) not null,
    created_at timestamp with time zone default NOW()
);

CREATE INDEX if not exists idx_orders_user_id ON orders(user_id);
