CREATE TABLE IF NOT EXISTS users (
    email text primary key,
    name text not null,
    picture_path text not null
);