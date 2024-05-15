create table devices (
    id bigint primary key generated always as identity,
    name text not null,
    token text not null,
    last_sync timestamp,
    sleeps_until timestamp
);

create table images (
    id bigint primary key generated always as identity,
    device_id bigint references devices(id) not null,
    permanent bool not null,
    data_original text not null,
    data_processed text not null
);