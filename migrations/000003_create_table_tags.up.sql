CREATE TABLE IF NOT EXISTS tags (
    id bigint primary key generated always as identity,
    tag_text text not null,
    event_id bigint references events(id) on delete cascade
);