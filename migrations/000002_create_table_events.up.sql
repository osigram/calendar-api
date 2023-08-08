CREATE TABLE IF NOT EXISTS events (
    id bigint primary key generated always as identity,
    color text not null,
    name text not null,
    description text not null,
    time_of_start timestamp not null,
    time_of_finish timestamp not null,
    email text references users(email) on delete cascade
);