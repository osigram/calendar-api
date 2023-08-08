CREATE TABLE IF NOT EXISTS extension_data (
    id bigint primary key generated always as identity,
    additional_data text not null,
    email text references users(email) on delete cascade
);